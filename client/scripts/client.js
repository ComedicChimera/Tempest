import 'https';
import { readFileSync } from 'fs';
import 'sjcl';

class Client {
    constructor(host, password) {
        this.hostname = host;
        this.ca = [ readFileSync('server-cert.pem') ];

        this.login(password)
            .then((token) => {
                this.token = token;
            })
            .catch((e) => {
                throw Error(e);
            });
    }  

    login(password) {
        return new Promise((resolve, reject) => {
            let options = {
                hostname: this.host,
                port: 443,
                path: '/get-nonce',
                method: 'GET',
                ca: this.ca,
            };
    
            options.agent = new https.Agent(options);

            this._performRequest(options)
                .then((ndata) => {
                    if (ndata.status) {
                        const snonce = Number(ndata.nonce);

                        let array = new Uint32Array(1);
                        crypto.getRandomValues(array);

                        const cnonce = array[0];

                        const hash = sjcl.hash.sha256.hash(password + String(snonce) + String(cnonce));
                        let tkOptions = {
                            ...options,
                            method: 'POST',
                            path: '/login',
                            body: JSON.stringify({
                                "cnonce": cnonce,
                                "auth-hash": hash
                            })
                        };
                        tkOptions.agent = new https.Agent(tkOptions);

                        this._performRequest(tkOptions)
                            .then((tokData) => {
                                if (tokData.status) {
                                    resolve(tokData.token);
                                } else {
                                    reject(tokData.message);
                                }
                            })
                            .catch((e) => {
                                reject(e);
                            });
                    } else {
                        reject(data.message);
                    }
                })
                .catch((e) => {
                    reject(e);
                });
        });
    }

    makeAPICall(route, method, data) {
        let options = {
            hostname: this.host,
            ca: this.ca,
            port: 443,
            path: route,
            method: method,
            body: JSON.stringify(data),
            headers: {
                "Authorization": this.token,
            }
        };

        return new Promise((resolve, reject) => {
            this._performRequest(options)
                .then((data) => {
                    if (data.status) {
                        resolve(data);
                    } else {
                        reject(data.message);
                    }
                })
                .catch((e) => {
                    reject(e);
                });
        });
    }

    _performRequest(options) {
        return new Promise((resolve, reject) => {
            const req = https.request(options, (res) => {
                if (res.statusCode != 200) {
                    reject(`Received response with bad status code of ${res.statusCode}`);
                    return;
                }

                res.on('data', (data) => {
                    resolve(JSON.parse(data));
                });
            });

            req.on('error', (e) => {
                reject(String(e));
            });
        });
    }

    sendFile(path, file) {

    }

    getFile(path) {

    }
}