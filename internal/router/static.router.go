package router

func (r *Router) staticRouter() {
	public := r.router.Group("public/images")
	public.Static("profile", "./public/photo_profiles")
}
