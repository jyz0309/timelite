package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var defaultDashboardPath = "./config/dashboards"

var dashboardsConfig map[string]*Dashboard

func init() {
	if _, err := os.Stat(defaultDashboardPath); os.IsNotExist(err) {
		os.MkdirAll(defaultDashboardPath, 0755)
	}
	dashboardsConfig = make(map[string]*Dashboard, 0)
	dashboardsConfig = ListDashboards()
}

func ListDashboards() map[string]*Dashboard {
	if len(dashboardsConfig) > 0 {
		return dashboardsConfig
	}

	files, err := os.ReadDir(defaultDashboardPath)
	if err != nil {
		return nil
	}
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		dashboard, err := GetDashboard(filepath.Join(defaultDashboardPath, file.Name()))
		if err != nil {
			panic(fmt.Sprintf(defaultDashboardPath+"/"+file.Name(), err, len(files)))
			// TODO
		}
		dashboardsConfig[dashboard.Name] = dashboard
	}
	return dashboardsConfig
}

type Dashboard struct {
	filename string
	Name     string            `json:"name"`
	Panels   map[string]*Panel `json:"panels"`
}

func NewDashboard(name string) (*Dashboard, error) {
	config, err := os.Create(filepath.Join(defaultDashboardPath, name+".json"))
	if err != nil {
		return nil, err
	}

	dashboard := &Dashboard{
		filename: config.Name(),
		Name:     name,
		Panels:   make(map[string]*Panel, 0),
	}

	bytes, err := json.Marshal(dashboard)
	if err != nil {
		return nil, err
	}
	_, err = config.Write(bytes)
	if err != nil {
		return nil, err
	}
	dashboardsConfig[dashboard.Name] = dashboard
	return dashboard, nil
}

func GetDashboard(filepath string) (*Dashboard, error) {
	config, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	dashboard := &Dashboard{
		filename: filepath,
		Panels:   make(map[string]*Panel, 0),
	}
	err = json.Unmarshal(config, dashboard)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

func DeleteDashboard(name string) error {
	dashboards := ListDashboards()
	for _, dashboard := range dashboards {
		if name != dashboard.Name {
			continue
		} else {
			if err := dashboard.Backup(); err != nil {
				return err
			}
			delete(dashboardsConfig, name)
			return os.Remove(dashboard.filename)
		}
	}
	return fmt.Errorf("failed to delete dashboard, err[%s]", "dashboard not exist")
}

func (b *Dashboard) AddPanel(panel *Panel) error {
	for _, p := range b.Panels {
		if p.Title == panel.Title {
			return fmt.Errorf("panel title already exists")
		}
	}
	b.Panels[panel.Title] = panel
	return nil
}

func (b *Dashboard) Save() error {
	bytes, err := json.Marshal(b)
	if err != nil {
		return err
	}
	if err := b.Backup(); err != nil {
		return err
	}
	return os.WriteFile(b.filename, bytes, 0644)
}

func (b *Dashboard) Backup() error {
	bytes, err := os.ReadFile(b.filename)
	if err != nil {
		return err
	}
	return os.WriteFile(b.filename+".bak", bytes, 0644)
}

func (b *Dashboard) GetAllPanels() map[string]*Panel {
	return b.Panels
}

func (b *Dashboard) DeletePanel(panel *Panel) error {
	if _, exist := b.Panels[panel.Title]; !exist {
		return fmt.Errorf("failed to delete panel, err[%s %s]", panel.Title, "not exist")
	}
	delete(b.Panels, panel.Title)
	return b.Save()
}

func (b *Dashboard) GetPanel(title string) *Panel {
	return b.Panels[title]
}

func (b *Dashboard) UpdatePanel(original *Panel, updated *Panel) error {
	if _, exist := b.Panels[original.Title]; !exist {
		return fmt.Errorf("failed to update panel, err[%s]", "panel not exist")
	}
	delete(b.Panels, original.Title)
	b.Panels[updated.Title] = updated
	return b.Save()
}

type Panel struct {
	Title     string `json:"title"`
	Query     string `json:"query"`
	PanelType string `json:"panel_type"`
}
