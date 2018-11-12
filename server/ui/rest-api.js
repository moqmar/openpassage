function RestApi(endpoint, options) {
  this.endpoint = (endpoint || "/").replace(/\/$/, "");
  this.options = options || {};
  // Allow usage in queue, or in other places where the value of "this" is incorrect:
  for (let method in ["GET", "POST", "PUT", "DELETE", "HEAD", "PATCH"].reduce(
    (p, c) => (p[c] = 1),
    {}
  )) {
    this[method] = RestApi.prototype[method].bind(this);
  }
}
RestApi.fetchResponseHandler = function fetchResponseHandler(hardFail, res) {
  return new Promise((resolve, reject) => {
    (res.headers.get("Content-Type").match(/^application\/json(?:$|;)/)
      ? res.json()
      : res.text()
    ).then(body => {
      res.content = body;
      if (!res.ok && hardFail) reject(res);
      else resolve(res);
    }, reject);
  });
};
RestApi.prototype.send = function send(method, path, headers, options) {
  return fetch(this.endpoint + path.replace(/^([^/])/, "/$1"), {
    ...this.options,
    ...options,
    method,
    headers: {
      ...this.options.headers,
      ...(options ? options.headers : {}),
      ...headers
    }
  }).then(
    RestApi.fetchResponseHandler.bind(
      null,
      this.options.hardFail
        ? (options || {}).hardFail !== false
        : (options || {}).hardFail
    )
  );
};
RestApi.prototype.sendWithBody = function sendWithBody(
  method,
  path,
  body,
  headers,
  options
) {
  let type = "text/plain";
  if (body instanceof FormData) type = "multipart/form-data";
  else if (body instanceof URLSearchParams)
    type = "application/x-www-form-urlencoded";
  else if (body instanceof Object) {
    type = "application/json";
    body = JSON.stringify(body);
  }

  return this.send(method, path, headers, {
    ...options,
    body,
    headers: { "Content-Type": type, ...(options ? options.headers : {}) }
  });
};
/**
 * Send a GET request to the API
 * @param {string} path Path of the REST endpoint to request
 * @param {object} headers Headers to set on the request
 * @param {object} options Other options to forward to fetch()
 * @returns {Promise<{body, response}>}
 */
RestApi.prototype.GET = function GET(path, headers, options) {
  return this.send("GET", path, headers, options);
};

/**
 * Send a POST request to the API
 * @param {string} path Path of the REST endpoint to request
 * @param {string|object|FormData|URLSearchParams} body Request body to send to the API; the content type will be determined automatically (text/plain, application/json, multipart/form-data, application/x-www-form-urlencoded)
 * @param {object} headers Headers to set on the request
 * @param {object} options Other options to forward to fetch()
 * @returns {Promise<{body, response}>}
 */
RestApi.prototype.POST = function POST(path, body, headers, options) {
  return this.sendWithBody("POST", path, body, headers, options);
};

/**
 * Send a PUT request to the API
 * @param {string} path Path of the REST endpoint to request
 * @param {string|object|FormData|URLSearchParams} body Request body to send to the API; the content type will be determined automatically (text/plain, application/json, multipart/form-data, application/x-www-form-urlencoded)
 * @param {object} headers Headers to set on the request
 * @param {object} options Other options to forward to fetch()
 * @returns {Promise<{body, response}>}
 */
RestApi.prototype.PUT = function PUT(path, body, headers, options) {
  return this.sendWithBody("PUT", path, body, headers, options);
};

/**
 * Send a DELETE request to the API
 * @param {string} path Path of the REST endpoint to request
 * @param {object} headers Headers to set on the request
 * @param {object} options Other options to forward to fetch()
 * @returns {Promise<{body, response}>}
 */
RestApi.prototype.DELETE = function DELETE(path, headers, options) {
  return this.send("DELETE", path, headers, options);
};

/**
 * Send a HEAD request to the API
 * @param {string} path Path of the REST endpoint to request
 * @param {object} headers Headers to set on the request
 * @param {object} options Other options to forward to fetch()
 * @returns {Promise<Response>}
 */
RestApi.prototype.HEAD = function HEAD(path, headers, options) {
  return this.send("HEAD", path, headers, options);
};

/**
 * Send a PATCH request to the API
 * @param {string} path Path of the REST endpoint to request
 * @param {string|object|FormData|URLSearchParams} body Request body to send to the API; the content type will be determined automatically (text/plain, application/json, multipart/form-data, application/x-www-form-urlencoded)
 * @param {object} headers Headers to set on the request
 * @param {object} options Other options to forward to fetch()
 * @returns {Promise<{body, response}>}
 */
RestApi.prototype.PATCH = function PATCH(path, body, headers, options) {
  return this.sendWithBody("PATCH", path, body, headers, options);
};

module.exports = RestApi;
