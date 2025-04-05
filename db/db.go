package db

//import "github.com/google/uuid"

//"Dashboarder/config"
//"Dashboarder/mongodb"
//"Dashboarder/siridb"
//"Dashboarder/valkeydb"

// our CLIs are primar connection handler for specific db type.
// we achieve them by getting them from db type api package constructor
type Databases struct {
	MongoCLI  interface{}
	SiriCLI   interface{}
	ValkeyCLI interface{}
}

type DBInterface interface {
	Create(string,interface{}) (interface{},error)
	Read(string) (interface{},error)
	Update(string, interface{}) (bool, error)
	Delete(key string) (bool, error)
}

/* there is 1 detail here.
 * We don't delete data.
 * We just mark them as deleted,so they won't be available directly.
 * in case we need to destroy some data, admin must do that manually!
 */

// Methods
// NewDBS -> Constructor for our builder
func NewDBS() *Databases {
	return new(Databases)
}

// Create(key string, v ...interface{}) (id interface{}, err error)
// Save data to database (create record in database) based on data type
func (dbs *Databases) Create(key string, v interface{}) (id interface{}, err error) {

	return
}

// Read(key string) (value interface{}, err error)
func (dbs *Databases) Read(key string) (value interface{}, err error) {

	return
}

// Update() -> no func now, need to delete/create ;)
func (dbs *Databases) Update(key string, newval interface{}) (updated bool, err error) {

	return
}

// Delete -> doing only setting up deleted flag, we dont delete data ;)
func (dbs *Databases) Delete(key string) (ok bool, err error) {

	return
}
