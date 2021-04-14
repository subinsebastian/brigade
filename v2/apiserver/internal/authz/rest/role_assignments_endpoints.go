package rest

import (
	"net/http"

	"github.com/brigadecore/brigade/v2/apiserver/internal/authz"
	libAuthz "github.com/brigadecore/brigade/v2/apiserver/internal/lib/authz"
	"github.com/brigadecore/brigade/v2/apiserver/internal/lib/restmachinery"
	"github.com/brigadecore/brigade/v2/apiserver/internal/system"
	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

// RoleAssignmentsEndpoints implements restmachinery.Endpoints to provide
// RoleAssignments-related URL --> action mappings to a restmachinery.Server.
type RoleAssignmentsEndpoints struct {
	AuthFilter                 restmachinery.Filter
	RoleAssignmentSchemaLoader gojsonschema.JSONLoader
	Service                    authz.RoleAssignmentsService
}

func (r *RoleAssignmentsEndpoints) Register(router *mux.Router) {
	// Grant a Role to a User or Service Account
	router.HandleFunc(
		"/v2/role-assignments",
		r.AuthFilter.Decorate(r.grant),
	).Methods(http.MethodPost)

	// Revoke a Role for a User or Service Account
	router.HandleFunc(
		"/v2/role-assignments",
		r.AuthFilter.Decorate(r.revoke),
	).Methods(http.MethodDelete)
}

func (r *RoleAssignmentsEndpoints) grant(
	w http.ResponseWriter,
	req *http.Request,
) {
	roleAssignment := libAuthz.RoleAssignment{}
	restmachinery.ServeRequest(
		restmachinery.InboundRequest{
			W:                   w,
			R:                   req,
			ReqBodySchemaLoader: r.RoleAssignmentSchemaLoader,
			ReqBodyObj:          &roleAssignment,
			EndpointLogic: func() (interface{}, error) {
				return nil, r.Service.Grant(req.Context(), roleAssignment)
			},
			SuccessCode: http.StatusOK,
		},
	)
}

func (r *RoleAssignmentsEndpoints) revoke(
	w http.ResponseWriter,
	req *http.Request,
) {
	roleAssignment := libAuthz.RoleAssignment{
		Role: libAuthz.Role{
			Type:  system.RoleTypeSystem,
			Name:  libAuthz.RoleName(req.URL.Query().Get("roleName")),
			Scope: req.URL.Query().Get("roleScope"),
		},
		Principal: libAuthz.PrincipalReference{
			Type: libAuthz.PrincipalType(req.URL.Query().Get("principalType")),
			ID:   req.URL.Query().Get("principalID"),
		},
	}
	restmachinery.ServeRequest(
		restmachinery.InboundRequest{
			W: w,
			R: req,
			EndpointLogic: func() (interface{}, error) {
				return nil, r.Service.Revoke(req.Context(), roleAssignment)
			},
			SuccessCode: http.StatusOK,
		},
	)
}
