package main

import (
	"fmt"
	"github.com/csby/security/certificate"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"path/filepath"
	"strings"
)

type Frame struct {
	*walk.MainWindow

	CenterScreen bool

	cfg     *Config
	cfgPath string
	dbCa    *walk.DataBinder
	dbCrt   *walk.DataBinder

	firstBoundsChange bool
	toolBar           *walk.ToolBar
	mainComposite     *walk.Composite
	kindSelector      *walk.ComboBox
	crtFileName       *walk.LineEdit
	createButton      *walk.PushButton
	hostButton        *walk.PushButton
}

func (s *Frame) OnBoundsChanged() {
	if !s.firstBoundsChange {
		s.firstBoundsChange = true

		if s.CenterScreen {
			screenWidth := int(win.GetSystemMetrics(win.SM_CXSCREEN))
			screenHeight := int(win.GetSystemMetrics(win.SM_CYSCREEN))
			frameWidth := s.Size().Width
			frameHeight := s.Size().Height
			frameBound := walk.Rectangle{
				X:      (screenWidth - frameWidth) / 2,
				Y:      (screenHeight - frameHeight) / 2,
				Width:  frameWidth,
				Height: frameHeight,
			}
			s.SetBounds(frameBound)
		}
	}
}

func (s *Frame) OnKindChanged() {
	selIndex := s.kindSelector.CurrentIndex()
	if selIndex < 0 {
		return
	}

	kind := Kinds[selIndex]
	if kind.Name == "server" {
		s.hostButton.SetVisible(true)
	} else {
		s.hostButton.SetVisible(false)
	}

	fileName := s.crtFileName.Text()
	if fileName == "" {
		s.cfg.Crt.Name = kind.Name
		s.crtFileName.SetText(kind.Name)
	} else if fileName != kind.Name {
		for _, item := range Kinds {
			if fileName == item.Name {
				s.crtFileName.SetText(kind.Name)
				break
			}
		}
	}
}

func (s *Frame) OnCreateCrt() {
	s.toolBar.SetEnabled(false)
	s.mainComposite.SetEnabled(false)
	s.createButton.SetText("创建中...")

	go func() {
		err := s.createCrt()
		if err != nil {
			walk.MsgBox(&s.FormBase, "新建证书失败", err.Error(), walk.MsgBoxIconError)
		} else {
			folder := filepath.Join(s.cfg.Crt.RootFolder, s.cfg.Crt.SubFolder)
			pfxPath := filepath.Join(folder, fmt.Sprintf("%s.pfx", s.cfg.Crt.Name))
			s.cfg.SaveToFile(s.cfgPath)
			walk.MsgBox(&s.FormBase, "新建证书成功", pfxPath, walk.MsgBoxIconInformation)
		}

		s.createButton.SetText("创建")
		s.mainComposite.SetEnabled(true)
		s.toolBar.SetEnabled(true)
	}()

}

func (s *Frame) SaveConfig() error {
	err := s.dbCa.Submit()
	if err != nil {
		return err
	}
	err = s.dbCrt.Submit()
	if err != nil {
		return err
	}

	return s.cfg.SaveToFile(s.cfgPath)
}

func (s *Frame) ShowHosts() {
	dlg := &DlgHost{cfg: s.cfg}
	err := dlg.Init(&s.FormBase)
	if err != nil {
		fmt.Println(err)
		return
	}
	dlg.ShowModal()
}

func (s *Frame) ShowCa() {
	dlg := &DlgCa{cfg: s.cfg, cfgPath: s.cfgPath, dbCa: s.dbCa}
	err := dlg.Init(&s.FormBase)
	if err != nil {
		fmt.Println(err)
		return
	}
	dlg.ShowModal()
}

func (s *Frame) ShowCrl() {
	caCrt, caKey, err := s.verifyCa()
	if err != nil {
		walk.MsgBox(&s.FormBase, "CA无效", err.Error(), walk.MsgBoxIconError)
		return
	}

	dlg := &DlgCrl{cfg: s.cfg, cfgPath: s.cfgPath, caCrt: caCrt, caKey: caKey}
	err = dlg.Init(&s.FormBase)
	if err != nil {
		fmt.Println(err)
		return
	}
	dlg.ShowModal()
}

func (s *Frame) VerifyCa() {
	crt, _, err := s.verifyCa()
	if err != nil {
		walk.MsgBox(&s.FormBase, "验证失败", err.Error(), walk.MsgBoxIconError)
	} else {
		msg := &strings.Builder{}
		msg.WriteString(fmt.Sprintf("证书类型：%s\r\n", kindDisplayName(crt.Organization())))
		msg.WriteString(fmt.Sprintf("证书标识：%s\r\n", crt.OrganizationalUnit()))
		msg.WriteString(fmt.Sprintf("显示名称：%s\r\n", crt.CommonName()))
		msg.WriteString(fmt.Sprintf("有效期：%s 至 %s", crt.NotBefore().Format("2006-01-02"), crt.NotAfter().Format("2006-01-02")))
		walk.MsgBox(&s.FormBase, "验证成功", msg.String(), walk.MsgBoxIconInformation)
		s.cfg.SaveToFile(s.cfgPath)
	}
}

func (s *Frame) createCrt() error {
	err := s.dbCrt.Submit()
	if err != nil {
		return err
	}
	if s.cfg.Ca.CrtFile == "" {
		return fmt.Errorf("CA证书文件为空")
	}
	if s.cfg.Ca.KeyFile == "" {
		return fmt.Errorf("CA私钥文件为空")
	}
	if s.cfg.Crt.Organization == "" {
		return fmt.Errorf("证书类型为空")
	}
	if s.cfg.Crt.RootFolder == "" {
		return fmt.Errorf("输出根目录为空")
	}
	if s.cfg.Crt.Name == "" {
		return fmt.Errorf("文件名称为空")
	}
	if s.cfg.Crt.OrganizationalUnit == "" {
		return fmt.Errorf("证书标识为空")
	}

	caCrt := &certificate.Crt{}
	err = caCrt.FromFile(s.cfg.Ca.CrtFile)
	if err != nil {
		return fmt.Errorf("加载CA证书错误: %v", err)
	}
	caPrivate := &certificate.RSAPrivate{}
	err = caPrivate.FromFile(s.cfg.Ca.KeyFile, fmt.Sprint(s.cfg.Ca.Password().Get()))
	if err != nil {
		return fmt.Errorf("加载CA私钥错误: %v", err)
	}

	private := &certificate.RSAPrivate{}
	err = private.Create(2048)
	if err != nil {
		return err
	}
	public, err := private.Public()
	if err != nil {
		return err
	}

	crtTemplate := &certificate.CrtTemplate{
		Organization:       s.cfg.Crt.Organization,
		OrganizationalUnit: s.cfg.Crt.OrganizationalUnit,
		CommonName:         s.cfg.Crt.CommonName,
		Locality:           s.cfg.Crt.Locality,
		Province:           s.cfg.Crt.Province,
		StreetAddress:      s.cfg.Crt.StreetAddress,
		Hosts:              s.cfg.Crt.Hosts,
	}
	template, err := crtTemplate.Template()
	if err != nil {
		return err
	}
	crt := &certificate.CrtPfx{}
	err = crt.Create(template, caCrt.Certificate(), public, caPrivate)
	if err != nil {
		return err
	}
	folder := filepath.Join(s.cfg.Crt.RootFolder, s.cfg.Crt.SubFolder)
	pfxPath := filepath.Join(folder, fmt.Sprintf("%s.pfx", s.cfg.Crt.Name))
	pfxPassword := fmt.Sprint(s.cfg.Crt.Password().Get())
	err = crt.ToFile(pfxPath, caCrt, private, pfxPassword)
	if err != nil {
		return err
	}
	crtPath := filepath.Join(folder, fmt.Sprintf("%s.crt", s.cfg.Crt.Name))
	crt.Crt.ToFile(crtPath)

	keyPath := filepath.Join(folder, fmt.Sprintf("%s.key", s.cfg.Crt.Name))
	private.ToFile(keyPath, "")

	return nil
}

func (s *Frame) verifyCa() (*certificate.Crt, *certificate.RSAPrivate, error) {
	err := s.dbCa.Submit()
	if err != nil {
		return nil, nil, err
	}

	key := &certificate.RSAPrivate{}
	err = key.FromFile(s.cfg.Ca.KeyFile, fmt.Sprint(s.cfg.Ca.Password().Get()))
	if err != nil {
		return nil, nil, fmt.Errorf("私钥无效: %v", err)
	}
	crt := &certificate.Crt{}
	err = crt.FromFile(s.cfg.Ca.CrtFile)
	if err != nil {
		return nil, nil, fmt.Errorf("证书无效: %v", err)
	}

	return crt, key, nil
}
