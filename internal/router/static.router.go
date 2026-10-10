package router

func (r *Router) staticRouter() {
	public := r.router.Group("public/images")
	public.Use(r.middleware.Cors)
	public.Static("profile", "./public/photo_profiles")
}
