package main

import (
	"encoding/json"
	"log"

	"github.com/Gimulator/client-go"
)

type Controller struct {
	Name      string
	Type      string
	Namespace string

	cli     *client.Client
	watcher chan client.Object
}

func NewController(name, namespace, typ string) (*Controller, error) {
	watcher := make(chan client.Object)
	cli, err := client.NewClient(watcher)
	if err != nil {
		return &Controller{}, err
	}

	return &Controller{
		Name:      name,
		Type:      typ,
		Namespace: namespace,
		cli:       cli,
		watcher:   watcher,
	}, nil
}

func (c *Controller) Run() {
	worldFilter := client.Key{
		Type:      c.Type,
		Name:      c.Name,
		Namespace: c.Namespace,
	}

	objs, err := c.cli.Find(worldFilter)
	switch {
	case err != nil:
		log.Printf("An error in finding: %s", err)
		return
	case len(objs) > 1:
		log.Println("Number of World is more than one")
		return
	case len(objs) == 1:
		world := World{}
		if err := json.Unmarshal([]byte(objs[0].Value.(string)), &world); err != nil {
			log.Printf("Cannot Structing the world object, Error: %s", err)
			return
		}
		c.watchWorld(world)
	}

	if err = c.cli.Watch(worldFilter); err != nil {
		log.Printf("Cannot watch on the world, Error: %s", err)
	}

	go func() {
		for obj := range c.watcher {
			world := World{}
			if err := json.Unmarshal([]byte(obj.Value.(string)), &world); err != nil {
				log.Printf("Cannot Structing the world object, Error: %s", err)
				continue
			}
			c.watchWorld(world)
		}
	}()
}

func (c *Controller) watchWorld(world World) {
	d := worldDrawer{
		World:  world,
		width:  width(),
		height: height(),
	}
	render(d)
	disableEvent = world.Turn != playerName
}
