package sendler

// type buttonData struct {
// 	Pressed bool `json:"pressed"`
// }

// func (s *sendler) sendButton(id types.Entity, broadcaster func(EntityInfo) error) error {

// 	c, _ := s.GetComponent(id, "button")
// 	button := c.(*components.ButtonComponent)

// 	broadcaster(EntityInfo{
// 		ID:   id,
// 		Type: "button",
// 		Data: buttonData{
// 			Pressed: button.Pressed,
// 		},
// 	})

// 	return nil
// }
