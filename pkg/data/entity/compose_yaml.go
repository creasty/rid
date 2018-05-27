package entity

// ComposeYaml represents a config file of docker-compose
type ComposeYaml struct {
	Rid struct {
		// ProjectName is used for `docker-compose` in order to distinguish projects in other locations
		ProjectName string `yaml:"project_name" valid:"required"`

		// MainService is a service name in `docker-compose.yml`, in which container commands given to rid are executed.
		// Default is "app"
		MainService string `yaml:"main_service" default:"app"`
	} `yaml:"rid" valid:"required"`

	Raw Hash `yaml:"-"`
}
