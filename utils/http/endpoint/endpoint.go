package endpoint

const (
	Root    string = "/"
	Health  string = "/health"
	Route   string = "/route"
	Service string = "/service"
)

var All = []string{
	Root,
	Health,
	Route,
	Service + Health,
}
