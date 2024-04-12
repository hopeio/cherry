package initialize

func (gc *globalConfig) Register() {
	/*	if gc.ConfigCenter == nil {
			return
		}
		svcName := gc.BasicConfig.Module
		_, err := gc.ConfigCenter.GetService(svcName)
		serviceConfig := gc.GetServiceConfig()
		if err != nil {
			err = gc.ConfigCenter.CreateService(svcName, &nacos.Metadata{
				Domain: serviceConfig.Domain,
				Port:   serviceConfig.Port,
			})
			if err != nil {
				log.Fatal(err)
			}
		}
		err = gc.ConfigCenter.RegisterInstance(serviceConfig.Port, svcName)
		if err != nil {
			log.Fatal(err)
		}*/
}
