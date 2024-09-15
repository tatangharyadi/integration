package xendit

import (
	"net/http"

	pubsub "github.com/tatangharyadi/pos-common/common/pubsub"
	message "github.com/tatangharyadi/pos-common/protobuf/message"
)

func (h Handler) CallbackQrPayment(w http.ResponseWriter, r *http.Request) {
	var qrPayment message.QrPayment
	qrPayment.Token = "diZjYadBT0KKj1YRZQwQp7:APA91bEnmxuYzvCdvx8uqIr0AlFbHfhPOBAQzdvVmi3JUfzzGaKeo56JIiyMTKmYtJOfx-2KiGNKM0OhBAjrwusFzGlWb7QhDwOGrKGOdgA3n5zycjjUzAeojqamw8NYxwzWp0NUHOyE"

	if err := pubsub.PublishProtoMessages(h.Env.GCPProjectId, h.Env.QrPaymentTopic, &qrPayment); err != nil {
		h.Logger.Error().Err(err).Msg("Error publish message to Pub/Sub")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(body)
}
