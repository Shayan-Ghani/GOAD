package controller

import (
	"fmt"
)


type Item struct {
	Id string
	name string
	description string
	tag []string
}

func (i *Item) AddItem(name string, desc string, tag ...string) {
	i.SetName(name)
	i.SetDesc(desc)
	i.SetTag(tag...)
}

func (i *Item) ViewItem(id string) string {
	
	return i.Id	
}
func (i *Item) ViewItems() (error, string, string, string, []string) {
	// turn to json
	return nil, i.Id, i.name, i.description, i.tag
}
// func (i *Item) ViewItemsDone(id string) string {
// 	return i.GetItemsDone(id)
// }

// func (i *Item) ViewItemByTag(tag string) string {
// 	return i.GetItemByTag(tag)
// }

// func (i *Item) GetItemByTag(tag string) (error , string) {

// 	// if !TagExists(tag){
// 	// 	return fmt.Errorf("tag doesn't exist"), ""
// 	// }
// 	/// join to get the tag name and item
// 	// err, tag_id = sql.Conn("select * from tags where name == ?")
// 	// err, item = sql.Conn("Select * from items where tag_id == ?")
// }



func (i *Item) DeleteItem(id string) Item{
	// p := i.GetItem()
	p := i
	
	// if err := mysql.Exec("delete from items where id = ?"); err == nil {
	// 	fmt.Errorf(err)
	// }

	return *p
}


/// Getter and setters

func (i *Item) GetId() string {
	return i.Id
}
// func (i *Item) GetName() string {
// 	return i.name
// }
func (i *Item) GetTag(tag string) string {
	return fmt.Sprintf(tag)
}
func (i *Item) GetTags() []string {
	return i.tag
}
// func (i *Item) GetDesc() string {
// 	return i.description
// }


func (i *Item) SetName(name string) {
	i.name = name
}
func (i *Item) SetTag(tag ...string) {
	i.tag= tag 
}

func (i *Item) SetDesc(desc string) {
	i.description = desc
}