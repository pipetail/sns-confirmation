package handlers

import (
	"log"
	"net/http"

	"github.com/pipetail/sns_verification/pkg/config"
	h "github.com/pipetail/sns_verification/pkg/http"
	"github.com/pipetail/sns_verification/pkg/sns/subscription"
)

func DefaultHandler(cfg config.Application, client h.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the message body
		v, err := subscription.VerificationFromReadCloser(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// log some info
		log.Printf("processing %s", v.SubscribeURL)

		// validate the URL
		err = subscription.ValidateURL(v.SubscribeURL, cfg.AllowedRegions, cfg.AllowedAccounts)
		if err != nil {
			log.Printf("validation failed: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// process the message body
		err = subscription.Verify(client, v.SubscribeURL, 200)
		if err != nil {
			log.Printf("could not process verification: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
