package api

import (
	context "golang.org/x/net/context"

	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ellcrys/util"
	"github.com/ncodes/cocoon/core/api/api/proto"
	"github.com/ncodes/cocoon/core/types"
)

// GetLogs fetches logs
func (api *API) GetLogs(ctx context.Context, req *proto.GetLogsRequest) (*proto.Response, error) {

	var err error
	var claims jwt.MapClaims

	if claims, err = api.checkCtxAccessToken(ctx); err != nil {
		return nil, types.ErrInvalidOrExpiredToken
	}

	cocoon, err := api.getCocoon(ctx, req.CocoonID)
	if err != nil {
		return nil, err
	}

	loggedInIdentity := claims["identity"].(string)

	// Ensure the cocoon identity matches the logged in user
	if cocoon.IdentityID != loggedInIdentity {
		return nil, fmt.Errorf("Permission denied: You do not have permission to perform this operation")
	}

	messages, err := api.logProvider.Get(ctx, fmt.Sprintf("connector-%s", req.CocoonID), int(req.NumLines), req.Source)
	if err != nil {
		return nil, err
	}

	messagesBytes, _ := util.ToJSON(messages)

	return &proto.Response{
		Status: 200,
		Body:   messagesBytes,
	}, nil
}