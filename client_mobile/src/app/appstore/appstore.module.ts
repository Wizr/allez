import { IonicModule } from 'ionic-angular';
import { AppStorePage } from './../../pages/appstore/appstore';
import { NgModule } from '@angular/core';
@NgModule({
    declarations:[
        AppStorePage
    ],
    entryComponents:[
        AppStorePage
    ],
    imports:[
        IonicModule.forRoot(AppStorePage)
    ]
})
export class AppStoreModule{}