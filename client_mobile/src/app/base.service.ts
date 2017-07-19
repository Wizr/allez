import { Injectable } from '@angular/core';
import { Headers, Http, RequestOptions } from '@angular/http';
import 'rxjs/add/operator/toPromise';

export class BaseService {
    constructor(public http: Http) {
    }

    private getOptions(isAuthenticated: boolean): Promise<RequestOptions> {
        // if (isAuthenticated) {
        //     return this.getUserToken().then((userToken) => {
        //         return new RequestOptions({
        //             headers: new Headers({
        //                 'Content-Type': 'application/json',
        //                 //'Token': userToken.Token
        //             })
        //         });
        //     });
        // }
        // else {
            return new Promise<RequestOptions>((resolve) => {
                return resolve(new RequestOptions({
                    headers: new Headers({
                        'Content-Type': 'application/json',
                    })
                }))
            });
        //}
    }

    // public getUserToken(): Promise<UserToken> {
    //     let token = sessionStorage.getItem("Token");
    //     if (token) {
    //         let changePasswordOnNextLogin = sessionStorage.getItem("ChangePasswordOnNextLogin");
    //         let userToken = new UserToken();
    //         userToken.Token = token;
    //         userToken.ChangePasswordOnNextLogin = changePasswordOnNextLogin == "true";

    //         return new Promise<UserToken>((resolve) => resolve(userToken));
    //     } else {
    //         return this.accountDB.getLatestToken().then((dt) => {
    //             let userToken = UserToken.from(dt);
    //             if (userToken) {
    //                 return userToken;
    //             }
    //         });
    //     }
    // }

    public get(url: string, isAuthenticated: boolean, parent: any): Promise<Object> {
        return this.getOptions(isAuthenticated)
            .then((requestOptions) => this.http.get(url, requestOptions).toPromise());
    }

    public post(url: string, isAuthenticated: boolean, body, parent: any): Promise<Object> {
        return this.getOptions(isAuthenticated)
            .then((requestOptions) => this.http.post(url, body, requestOptions).toPromise());
    }

    public put(url: string, isAuthenticated: boolean, body, parent: any): Promise<Object> {
        return this.getOptions(isAuthenticated)
            .then((requestOptions) => this.http.put(url, body, requestOptions).toPromise());
    }

    // public uploadFile(url: string, fileUri: string, fileName, params?: { [s: string]: any }): Promise<Object> {
    //     return this.getUserToken().then((userToken) => {
    //         let fileTransfer = new Transfer().create();
    //         let options: FileUploadOptions = {
    //             fileKey: 'file',
    //             fileName: fileName,
    //             headers: {
    //                 'APIKey': GlobalVariables.API_KEY,
    //                 'FingerPrint': GlobalVariables.FINGERPRINT,
    //                 'Token': userToken.Token
    //             },
    //             params: params
    //         }

    //         return fileTransfer.upload(fileUri, url, options)
    //     });
    // }
}