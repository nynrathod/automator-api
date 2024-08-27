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
	"time"
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
		fmt.Println("current time: ", time.Now())
		// Access the OtpRequest and WebSocket instance
		otpRequest := requestWithWS.OtpRequest
		uuid := requestWithWS.UUID

		// Create context with UUID
		ctx := context.WithValue(context.Background(), "UUID", uuid)

		// Call ProcessOtpRequest and pass ws if needed
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
		fmt.Printf("Custom event - User: %s\n", ep.Kws.GetStringAttribute("user_id"))
		//SendNoti()
	})

	socketio.On("sms_receive", func(ep *socketio.EventPayload) {
		fmt.Printf("SMS receive event\n")
		//SendNoti()
	})

	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		fmt.Println("Message", string(ep.Data))
		// Access the WebSocket instance through the EventPayload
		ws := ep.Kws
		uuid := ws.UUID

		// Now you can access the UUID or other properties of the WebSocket instance
		//fmt.Println("WebSocket UUID:", ws.UUID)

		//message := fmt.Sprintf("Welcome! Your WebSocket UUID is: %s", ws.UUID)
		//ws.EmitTo(ws.UUID, []byte(message), socketio.TextMessage)

		// Determine event type and handle appropriately
		var genericMessage map[string]interface{}
		if err := json.Unmarshal(ep.Data, &genericMessage); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			return
		}

		eventType, ok := genericMessage["event"].(string)
		if !ok {
			fmt.Println("Event type missing or invalid")
			return
		}

		if err := handleEvent(eventType, ep.Data, db, uuid); err != nil {
			fmt.Println("Error handling event:", err)
		}

		fmt.Println("eventType", eventType)

		// Unmarshal the JSON payload into GenericMessage
		//var message *entities.Sms
		//err := json.Unmarshal(ep.Data, &message)
		//if err != nil {
		//	fmt.Println("Error unmarshalling message:", err)
		//	return
		//}
		//
		//switch message.Event {
		//case "SMS_RECEIVE":
		//	message.UserId = ep.Kws.GetStringAttribute("user_id")
		//	sms.StoreSms(message, db)
		//	//handleCustomEvent(message.Data)
		//case "REQUEST_OTP":
		//	message.UserId = ep.Kws.GetStringAttribute("user_id")
		//	sms.StoreSms(message, db)
		//default:
		//	//handleOtherEvent(message.Data)
		//}

		//fmt.Printf("Message event - User: %s - Message: %s\n", ep.Kws.GetStringAttribute("user_id"), string(ep.Data))
		//
		//message := MessageObject{}
		//err := json.Unmarshal(ep.Data, &message)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//
		//if message.Event != "" {
		//	ep.Kws.Fire(message.Event, []byte(message.Data))
		//}
		//
		//err = ep.Kws.EmitTo(clients[message.To], ep.Data, socketio.TextMessage)
		//if err != nil {
		//	fmt.Println(err)
		//}
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
		//kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", userId, kws.UUID)), true, socketio.TextMessage)
		//kws.EmitTo(kws.UUID, []byte("sdas"))
		//kws.Emit([]byte(fmt.Sprintf("Hello user: %s with UUID: %s", userId, kws.UUID)), socketio.TextMessage)
	}))
}

func handleEvent(eventType string, data []byte, db *mongo.Database, uuid string) error {
	switch eventType {
	case "SMS_RECEIVE":

		//err := kws.EmitTo("8a84d790-d3b9-4113-9df9-c3fa53f16293", []byte("sdas"))
		//fmt.Println("myrrer", err)
		var smsData *entities.Sms
		if err := json.Unmarshal(data, &smsData); err != nil {
			return fmt.Errorf("error unmarshalling SMS_RECEIVE: %v", err)
		}
		//sms.UserId = ep.Kws.GetStringAttribute("user_id")
		sms.StoreSms(smsData, db)
		fmt.Println("Handling SMS_RECEIVE:", smsData.UserId)
		// Perform operations like storing SMS
	case "REQUEST_OTP":

		//message := fmt.Sprintf("Welcome! Your WebSocket UUID is: %s", ws.UUID)
		//ws.EmitTo(ws.UUID, []byte(message), socketio.TextMessage)
		var otp *entities.OtpRequest
		if err := json.Unmarshal(data, &otp); err != nil {
			return fmt.Errorf("error unmarshalling REQUEST_OTP: %v", err)
		}
		otpRequestWithWS := &OtpRequestWithWS{
			OtpRequest: otp,
			UUID:       uuid,
		}

		// Push the request and WebSocket instance to the queue
		notificationQueue <- otpRequestWithWS
		//sms.ProcessOtpRequest(otp, db)
		fmt.Println("Handling REQUEST_OTP:", otp)
		// Perform operations like handling OTP request
	// Add more cases for other events
	default:
		return fmt.Errorf("unhandled event type: %s", eventType)
	}
	return nil
}
