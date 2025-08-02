package endpoint

const (
	Root    string = "/"
	Health  string = "/health"
	Route   string = "/route"
	Service string = "/service"
	Metrics string = "/metrics"
)

var All = []string{
	Root,
	Health,
	Route,
	Metrics,
	Service + Health,
	Service + Metrics,
}
