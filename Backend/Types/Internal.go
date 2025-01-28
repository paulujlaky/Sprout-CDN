package Types

type Config struct {

	GinMode string

	Server struct {

		Port int 

	}

	Database struct {

		URI string 
		Name string 

	}

	Logging struct {

		Mode string
		Webhook string

	}

	Versions struct {

		Server string
		Frontend string

	}

}

