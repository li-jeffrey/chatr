package controllers

import (
	"chatr/internal/store"

	"github.com/valyala/fasthttp"
)

func SubmitQuestion(ctx *fasthttp.RequestCtx) {
	sessionID := ctx.QueryArgs().Peek("sessionID")
	if sessionID == nil {
		ctx.Error("SessionID required", fasthttp.StatusBadRequest)
		return
	}

	store.CreateSubmission(ctx.PostBody(), sessionID)
}

func SubmitAnswer(ctx *fasthttp.RequestCtx) {
	id := string(ctx.QueryArgs().Peek("id"))
	if e := store.UpdateSubmission(id, ctx.PostBody()); e != nil {
		ctx.Error(e.Error(), fasthttp.StatusConflict)
	}
}
