function toggleNavigation() {
    $("#sidebar-toggler").on("click", _ => {
        $("nav").toggleClass("opened");
    });
}

window.toggleNavigation = toggleNavigation;
