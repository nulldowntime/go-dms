package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	auth "github.com/ninadakolekar/go-dms/src/auth"
	"github.com/ninadakolekar/go-dms/src/constants"
	"github.com/ninadakolekar/go-dms/src/docs"
)

// ProcessReviewApproval ... Handles Review Approval
func ProcessReviewApproval(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		// User Auth Verification

		username, err := auth.GetCurrentUser(r)

		if err != nil { // Auth unsucessful
			fmt.Println("ERROR ProcessReviewApproval Line 23: ", err) // Debug
			http.Redirect(w, r, "/", 302)
		}

		vars := mux.Vars(r)
		id := vars["id"]

		if docs.ValidateDocNo(id) {

			Document, err := docs.FetchDocByID(id)

			if err != nil {

				fmt.Println("Failed to fetch document: ", err)
				return
			}

			if !docs.CheckCurrentReviewer(Document, username) || Document.FlowStatus != constants.ReviewFlow {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			rwResponse := r.FormValue("rw-answer")

			if rwResponse == "approve" {

				if Document.DocProcess == constants.Everyone {

					Document.FlowList = append(Document.FlowList, username)

					if len(Document.FlowList) == len(Document.Reviewer) {
						Document.FlowStatus = constants.ApproveFlow
						Document.FlowList = nil
						Document.CurrentFlowUser = 0
					}

					_, err := docs.AddInactiveDoc(Document)

					if err != nil {
						fmt.Println("Error: Document Update Failed ProcessReviewApproval Line 60 ", err) // Debug
						return
					}

					fmt.Fprintf(w, "Approval Recorded.")

				} else if Document.DocProcess == constants.Anyone {

					Document.FlowStatus = constants.ApproveFlow
					Document.FlowList = nil
					Document.CurrentFlowUser = 0

					_, err := docs.AddInactiveDoc(Document)

					if err != nil {
						fmt.Println("Error: Document Update Failed ProcessReviewApproval Line 73 ", err) // Debug
						return
					}

					fmt.Fprintf(w, "Approval Recorded.")

				} else if Document.DocProcess == constants.OneByOne {

					Document.CurrentFlowUser++

					if int(Document.CurrentFlowUser) == len(Document.Reviewer) {
						Document.FlowStatus = constants.ApproveFlow
						Document.FlowList = nil
						Document.CurrentFlowUser = 0
					}

					_, err := docs.AddInactiveDoc(Document)

					if err != nil {
						fmt.Println("Error: Document Update Failed ProcessReviewApproval Line 92 ", err) // Debug
						return
					}

					fmt.Fprintf(w, "Approval Recorded.")

				}

			} else if rwResponse == "reject" {

				if Document.DocProcess == constants.Everyone {

					Document.FlowStatus = constants.CreateFlow
					Document.FlowList = nil
					Document.CurrentFlowUser = 0

					_, err := docs.AddInactiveDoc(Document)

					if err != nil {
						fmt.Println("Error: Document Update Failed ProcessReviewApproval Line 109 ", err) // Debug
						return
					}

					fmt.Fprintf(w, "Rejection Recorded.")

				} else if Document.DocProcess == constants.Anyone {

					Document.FlowList = append(Document.FlowList, username)

					if len(Document.FlowList) == len(Document.Reviewer) {
						Document.FlowStatus = constants.CreateFlow
						Document.FlowList = nil
						Document.CurrentFlowUser = 0
					}

					_, err := docs.AddInactiveDoc(Document)

					if err != nil {
						fmt.Println("Error: Document Update Failed ProcessReviewApproval Line 126 ", err) // Debug
						return
					}

					fmt.Fprintf(w, "Rejection Recorded.")

				} else if Document.DocProcess == constants.OneByOne {

					Document.FlowStatus = constants.CreateFlow
					Document.FlowList = nil
					Document.CurrentFlowUser = 0

					_, err := docs.AddInactiveDoc(Document)

					if err != nil {
						fmt.Println("Error: Document Update Failed ProcessReviewApproval Line 139 ", err) // Debug
						return
					}

					fmt.Fprintf(w, "Rejection Recorded.")

				}

			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}

		}

	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
