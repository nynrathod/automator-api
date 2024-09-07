package bootstrap

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/nynrathod/automator-api/api/handlers/sms"
	"github.com/nynrathod/automator-api/pkg/entities"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

var clients = make(map[string]string)

type OtpRequestWithWS struct {
	OtpRequest *entities.OtpRequest
	UUID       string
}

// Change notificationQueue to hold OtpRequestWithWS
var notificationQueue = make(chan *OtpRequestWithWS, 100)

func ProcessQueue(db *mongo.Database) {
	for requestWithWS := range notificationQueue {

		otpRequest := requestWithWS.OtpRequest
		uuid := requestWithWS.UUID

		ctx := context.WithValue(context.Background(), "UUID", uuid)
		_, err := sms.ProcessOtpRequest(otpRequest, db, ctx)
		if err != nil {
			fmt.Printf("Error processing OTP: %v\n", err)
		}

	}
}

func SetupWebSocket(app *fiber.App, db *mongo.Database) {
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Printf("Connection event - User: %s\n", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On("CUSTOM_EVENT", func(ep *socketio.EventPayload) {
		//fmt.Printf("Custom event - User: %s\n", ep.Kws.GetStringAttribute("user_id"))
		//SendNoti()
	})

	socketio.On("sms_receive", func(ep *socketio.EventPayload) {
		//fmt.Printf("SMS receive event\n")
		//SendNoti()
	})

	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {

		ws := ep.Kws
		uuid := ws.UUID
		// Now you can access the UUID or other properties of the WebSocket instance
		//fmt.Println("WebSocket UUID:", ws.UUID)

		//message := fmt.Sprintf("Welcome! Your WebSocket UUID is: %s", ws.UUID)
		//ws.EmitTo(ws.UUID, []byte(message), socketio.TextMessage)

		// Determine event type and handle appropriately
		var genericMessage map[string]interface{}
		if err := json.Unmarshal(ep.Data, &genericMessage); err != nil {
			//fmt.Println("Error unmarshalling message:", err)
			return
		}

		eventType, ok := genericMessage["event"].(string)
		if !ok {
			//fmt.Println("Event type missing or invalid")
			return
		}

		if err := handleEvent(eventType, ep.Data, db, uuid); err != nil {
			fmt.Println("Error handling event:", err)
		}

		//fmt.Println("eventType", eventType)
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		delete(clients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Disconnection event - User: %s\n", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		delete(clients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Close event - User: %s\n", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On(socketio.EventError, func(ep *socketio.EventPayload) {
		fmt.Printf("Error event - User: %s\n", ep.Kws.GetStringAttribute("user_id"))
	})

	app.Get("/ws/:id", socketio.New(func(kws *socketio.Websocket) {
		userId := kws.Params("id")
		clients[userId] = kws.UUID

		kws.SetAttribute("user_id", userId)
		fmt.Println("connecteuserId", kws.UUID)
	}))
}

func handleEvent(eventType string, data []byte, db *mongo.Database, uuid string) error {
	switch eventType {
	case "SMS_RECEIVE":

		var smsData *entities.Sms
		if err := json.Unmarshal(data, &smsData); err != nil {
			return fmt.Errorf("error unmarshalling SMS_RECEIVE: %v", err)
		}
		sms.StoreSms(smsData, db, uuid)
		//fmt.Println("Handling SMS_RECEIVE:", smsData.UserId)
		// Perform operations like storing SMS
	case "REQUEST_OTP":

		var otp *entities.OtpRequest
		if err := json.Unmarshal(data, &otp); err != nil {
			return fmt.Errorf("error unmarshalling REQUEST_OTP: %v", err)
		}
		otpRequestWithWS := &OtpRequestWithWS{
			OtpRequest: otp,
			UUID:       uuid,
		}

		notificationQueue <- otpRequestWithWS

		fmt.Println("Handling REQUEST_OTP:", otp)

	default:
		return fmt.Errorf("unhandled event type: %s", eventType)
	}
	return nil
}
