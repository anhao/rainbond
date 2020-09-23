package handler

import (
	"github.com/goodrain/rainbond/api/model"
	"github.com/goodrain/rainbond/api/util/bcode"
	"github.com/goodrain/rainbond/db"
	dbmodel "github.com/goodrain/rainbond/db/model"
	"github.com/sirupsen/logrus"
)

// AddConfigGroup -
func (a *ApplicationAction) AddConfigGroup(appID string, req *model.ApplicationConfigGroup) (*model.ApplicationConfigGroupResp, error) {
	var serviceResp []dbmodel.ServiceConfigGroup
	// Create application configGroup-services
	for _, sid := range req.ServiceIDs {
		services, _ := db.GetManager().TenantServiceDao().GetServiceByID(sid)
		serviceConfigGroup := dbmodel.ServiceConfigGroup{
			AppID:           appID,
			ConfigGroupName: req.ConfigGroupName,
			ServiceID:       sid,
			ServiceAlias:    services.ServiceAlias,
		}
		serviceResp = append(serviceResp, serviceConfigGroup)
		if err := db.GetManager().ServiceConfigGroupDao().AddModel(&serviceConfigGroup); err != nil {
			if err == bcode.ErrServiceConfigGroupExist {
				logrus.Warningf("config group \"%s\" under this service \"%s\" already exists.", serviceConfigGroup.ConfigGroupName, serviceConfigGroup.ServiceID)
				continue
			}
			return nil, err
		}
	}

	// Create application configGroup-configItem
	for _, it := range req.ConfigItems {
		configItem := &dbmodel.ConfigItem{
			AppID:           appID,
			ConfigGroupName: req.ConfigGroupName,
			ItemKey:         it.ItemKey,
			ItemValue:       it.ItemValue,
		}
		if err := db.GetManager().ConfigItemDao().AddModel(configItem); err != nil {
			if err == bcode.ErrConfigItemExist {
				logrus.Warningf("config item \"%s\" under this config group \"%s\" already exists.", configItem.ItemKey, configItem.ConfigGroupName)
				continue
			}
			return nil, err
		}
	}

	// Create application configGroup
	config := &dbmodel.ApplicationConfigGroup{
		AppID:           appID,
		ConfigGroupName: req.ConfigGroupName,
		DeployType:      req.DeployType,
	}
	if err := db.GetManager().ApplicationConfigDao().AddModel(config); err != nil {
		return nil, err
	}

	appconfig, _ := db.GetManager().ApplicationConfigDao().GetConfigByID(appID, req.ConfigGroupName)
	var resp *model.ApplicationConfigGroupResp
	resp = &model.ApplicationConfigGroupResp{
		CreateTime:      appconfig.CreatedAt,
		AppID:           appID,
		ConfigGroupName: appconfig.ConfigGroupName,
		DeployType:      appconfig.DeployType,
		ConfigItems:     req.ConfigItems,
		Services:        serviceResp,
	}
	return resp, nil
}
