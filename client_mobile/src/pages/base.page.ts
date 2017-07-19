import { AlertController, NavController, LoadingController, Content, Header, Loading } from 'ionic-angular';
import { ViewChild } from '@angular/core';

export class BasePage {
    @ViewChild(Content) content: Content;
    @ViewChild(Header) header: Header;
    public loader: Loading;

    private runPromise(fun: any): Promise<any> {
        return new Promise<any>((resolve) => resolve(fun(this)));
    }

    // private logoutWithError(error: any, title?: string): void {
    //     let isLogout: boolean = false;
    //     let alert = this.alertCtrl.create({
    //         title: title ? "Oops!" : title,
    //         message: error.toString(),
    //         enableBackdropDismiss: false,
    //         buttons: [{
    //             text: 'OK',
    //             handler: () => {
    //                 if (!isLogout) {
    //                     isLogout = true;
    //                     this.logout();
    //                 }
    //             }
    //         }]
    //     });
    //     alert.present();
    // }

    private alertConfirm(message: string, title?: string, fun?: any): void {
        let alert = this.alertCtrl.create({
            title: title ? "Confirm" : title,
            message: message.toString(),
            enableBackdropDismiss: false,
            buttons: [
                {
                    text: 'Cancel',
                    role: 'cancel'
                },
                {
                    text: 'OK',
                    handler: () => {
                        if (fun) {
                            fun();
                        }
                    }
                }]
        });
        alert.present();
    }



    constructor(public alertCtrl: AlertController, public navCtrl: NavController, public loadingCtrl: LoadingController) {
    }

    // ionViewDidEnter(): void {
    //     if (this.header) {
    //         // ionic margin top is wrong
    //         this.content.getNativeElement().querySelector(".scroll-content").style.marginTop = this.header.height() + 'px';
    //         this.content.getNativeElement().querySelector(".fixed-content").style.marginTop = this.header.height() + 'px';
    //     }
    // }

    public execute(fun: any): void {
        this.loader = this.loadingCtrl.create({
            content: "<div class='custom-spinner - container'><div class='custom-spinner-box' > 起飞中...</div></div>",
            dismissOnPageChange: true,
            showBackdrop: true
        });
        this.loader.present();

        this.runPromise(fun);
    }

    public dismissLoader() {
        this.loader.dismiss();
    }

    public refresh(refresher, fun: any): void {
        setTimeout(() => {
            this.runPromise(fun).then(() => {
                refresher.complete();
            });
        }, 1000);
    }

    // public logout(): void {
    //     this.execute(function (parent) {
    //         parent.accountService.logout()
    //             .then(response => parent.processResponse(response))
    //             .then(() => {
    //                 parent.dismissLoader()
    //                 parent.localLogout(LoginPage);
    //             });
    //     });
    // }

    // public localLogout(page?: any): void {
    //     sessionStorage.clear();
    //     this.accountDB.logout();

    //     if (page) {
    //         this.navCtrl.setRoot(page, {}, { animate: true, direction: "forward" });
    //     }
    // }

    public alertError(error: any, title?: string): void {
        let alert = this.alertCtrl.create({
            title: title ? "Oops!" : title,
            message: error.toString(),
            enableBackdropDismiss: false,
            buttons: [{
                text: 'OK',
            }]
        });
        alert.present();
    }

    public alertMessage(message: string, title?: string, fun?: any): void {
        let alert = this.alertCtrl.create({
            title: title ? "" : title,
            message: message.toString(),
            enableBackdropDismiss: false,
            buttons: [{
                text: 'OK',
                handler: () => {
                    if (fun) {
                        fun();
                    }
                }
            }]
        });
        alert.present();
    }

    public alertDelete(message?: string, title?: string, fun?: any): void {
        this.alertConfirm(message ? message : 'Are you sure you want to delete?', title, fun);
    }

    // public processAuthenticatedResponse(response, reject?: Function): Object {
    //     let appResponse = response.json() as AppResponse;
    //     if (!appResponse.IsAuthenticated && !appResponse.IsSuccessful) {
    //         if (reject) {
    //             reject();
    //         }

    //         this.logoutWithError(appResponse.Message, appResponse.Title);
    //     }
    //     else if (!appResponse.IsSuccessful) {
    //         if (reject) {
    //             reject();
    //         }

    //         this.alertError(appResponse.Message, appResponse.Title);
    //     }
    //     else {
    //         return appResponse.Result;
    //     }
    // }

    // public processResponse(response): Object {
    //     let appResponse = response.json() as AppResponse;
    //     if (!appResponse.IsSuccessful) {
    //         this.alertError(appResponse.Message, appResponse.Title);
    //     }
    //     else {
    //         return appResponse.Result;
    //     }
    // }

    public handleError(error: any, title: string) {
        if (error && error.status != 200) {
            this.alertError("The network connection to your application was unsuccessful. Please try again later.", "Network connection error");
        }
        else {
            this.alertError(error, title);
        }
    }
}

