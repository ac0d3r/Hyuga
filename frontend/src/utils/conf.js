const apihost = function () {
    const host = window.location.href;
    return host.endsWith("/") ? host.substring(0, host.length - 1) : host;
}()

export { apihost };
