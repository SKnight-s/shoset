package shoset

import (
	"fmt"
	"time"

	"github.com/ditrit/shoset/msg"
)

// GetConfigLink :
func GetConfigLink(c *ShosetConn) (msg.Message, error) {
	var cfg msg.ConfigLink
	err := c.ReadMessage(&cfg)
	return cfg, err
}

// HandleConfigLink :
func HandleConfigLink(c *ShosetConn, message msg.Message) error {
	// fmt.Println(c.ch.GetBindAddress(), " enter handleconfiglink !!!!!!!!!")
	cfg := message.(msg.ConfigLink)
	remoteAddress := cfg.GetAddress()
	dir := c.GetDir()
	switch cfg.GetCommandName() {
	case "link":
		if dir == "in" { // a socket wants to link to this one
			if connsJoin := c.ch.ConnsByName.Get(c.ch.GetLogicalName()); connsJoin != nil { //already linked
				if connsJoin.Get(remoteAddress) != nil {
					return nil
				}
			}

			c.SetRemoteAddress(remoteAddress)                            // avoid tcp port name
			c.ch.ConnsByName.Set(cfg.GetLogicalName(), remoteAddress, c) // set conn in this socket
			c.SetRemoteLogicalName(cfg.GetLogicalName())

			// fmt.Println(c.ch)
			// fmt.Println(c.ch.ConnsByName)

			localBrothers := c.ch.ConnsByName.Get(c.ch.GetLogicalName())
			// fmt.Println("local brothers : ", localBrothers)
			localBrothersArray := []string{}
			if localBrothers != nil {
				time.Sleep(time.Millisecond * time.Duration(1))
				localBrothersArray = localBrothers.Keys("all")
			}

			remoteBrothers := c.ch.ConnsByName.Get(cfg.GetLogicalName())
			remoteBrothersArray := []string{}
			if remoteBrothers != nil {
				// fmt.Println("~~~~~~~~")
				// fmt.Println(remoteBrothers)
				time.Sleep(time.Millisecond * time.Duration(1))
				remoteBrothersArray = remoteBrothers.Keys("all")
				// fmt.Println("~~~~~~~~")
			}
			fmt.Println(c.ch.GetBindAddress(), "remote brothers 1 : ", remoteBrothersArray)

			// fmt.Println("brothers arrays : ", localBrothersArray, remoteBrothersArray)

			brothers := msg.NewCfgBrothers(localBrothersArray, remoteBrothersArray, c.ch.GetLogicalName(), "brothers")
			c.SendMessage(brothers)
		}

	case "brothers":
		// fmt.Println(c.ch.GetBindAddress(), " enter case brothers")
		if dir == "out" { // this socket wants to link to another
			// fmt.Println("config name : ", cfg.GetLogicalName())
			c.ch.ConnsByName.Set(cfg.GetLogicalName(), c.GetRemoteAddress(), c) // set conns in the other socket
			// c.ch.ConnsByName.Set(c.ch.GetLogicalName(), c.GetRemoteAddress(), c) // set conns in the other socket
			c.SetRemoteLogicalName(cfg.GetLogicalName())

			// fmt.Println(c.ch)
			// fmt.Println(c.ch.ConnsByName)

			localBrothers := cfg.GetYourBrothers()
			// fmt.Println("local brothers 2: ", localBrothers)
			// fmt.Println("I'm ", c.ch.GetBindAddress(), " and here are my known local brothers : ", myKnownLocalBrothers)
			// fmt.Println(c.ch.ConnsByName)
			// fmt.Println(c.ch.ConnsByName.Get(c.ch.GetLogicalName()))
			remoteBrothers := cfg.GetMyBrothers()
			// fmt.Println("remote brothers 2 : ", remoteBrothers)
			// myKnownRemoteBrothers := c.ch.ConnsByName.Get(cfg.GetLogicalName())

			for _, bro := range localBrothers { // à tester en rajoutant un aga ////////////////// ne fonctionne pas encore
				if bro != c.ch.GetBindAddress() {
					conn, err := NewShosetConn(c.ch, bro, "me") // create empty socket so that the two aga know each other
					if err == nil {
						// fmt.Println(c.ch.GetBindAddress(), " has a new bro : ", bro, "####################")
						c.ch.ConnsByName.Set(c.ch.GetLogicalName(), bro, conn) // put them into ConnsByName - need to put this one in the other socket
						// send aknoledge_brother ???
					}
					newLocalBrothers := c.ch.ConnsByName.Get(c.ch.GetLogicalName()).Keys("me")
					for _, lName := range c.ch.ConnsByName.Keys() {
						lNameConns := c.ch.ConnsByName.Get(lName)
						addresses := lNameConns.Keys("in")
						brothers := msg.NewCfgBrothers(newLocalBrothers, addresses, c.ch.GetLogicalName(), "brothers")
						lNameConns.Iterate(
							func(key string, val *ShosetConn) {
								val.SendMessage(brothers)
							})
					}
				}
			}

			for _, remoteBro := range remoteBrothers { // link to the brothers of the socket it's linked with
				remoteBrothers := c.ch.ConnsByName.Get(cfg.GetLogicalName())
				if remoteBrothers != nil {
					if remoteBrothers.Get(remoteBro) == nil {
						// fmt.Println("!!!!!!!!!!! new link")
						c.ch.Link(remoteBro)
					}
				}
			}
		}
	}
	return nil
}
