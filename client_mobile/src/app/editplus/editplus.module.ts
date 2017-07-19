import { IonicModule } from 'ionic-angular';
import { EditPlusPage } from './../../pages/editplus/editplus';
import { NgModule } from '@angular/core';
@NgModule({
    declarations:[
        EditPlusPage
    ],
    entryComponents:[
        EditPlusPage
    ],
    imports:[
        IonicModule.forRoot(EditPlusPage)
    ]
})
export class EditPlusModule{}