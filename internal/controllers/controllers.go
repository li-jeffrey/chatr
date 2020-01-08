package controllers

import (
	"chatr/internal/store"

	"github.com/valyala/fasthttp"
)

func SubmitQuestion(ctx *fasthttp.RequestCtx) {
	s := store.CreateSubmission(ctx.PostBody())
	id, _ := s.ID.MarshalText()
	ctx.Response.SetBodyRaw(id)
}

func SubmitAnswer(ctx *fasthttp.RequestCtx) {
	id := string(ctx.QueryArgs().Peek("id"))
	if e := store.UpdateSubmission(id, ctx.PostBody()); e != nil {
		ctx.Response.SetBodyString(e.Error())
	}
}
