export class AuthRedirectError extends Error {
    constructor() {
        super("auth redirect");
        this.name = "AuthRedirectError";
    }
}

let reloading = false;

export async function apiFetch(input: RequestInfo | URL, init: RequestInit = {}): Promise<Response> {
    const res = await fetch(input, {
        credentials: "same-origin",
        ...init,
        redirect: "manual",
    });

    // Same-origin fetch + 3xx → opaqueredirect (status 0). With Caddy + tinyauth
    // the 302 points to a different origin, so we can't follow it from JS — we
    // need a top-level navigation so the browser follows the redirect normally
    // and the user lands on the auth page.
    if (res.type === "opaqueredirect") {
        if (!reloading) {
            reloading = true;
            window.location.assign(window.location.href);
        }
        throw new AuthRedirectError();
    }

    return res;
}
