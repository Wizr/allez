import { CharlesService } from './charles.service';
import { IonicModule } from 'ionic-angular';
import { CharlesPage } from './../../pages/charles/charles';
import { NgModule } from '@angular/core';
@NgModule({
    declarations:[
        CharlesPage
    ],
    entryComponents:[
        CharlesPage
    ],
    imports:[
        IonicModule.forRoot(CharlesPage)
    ],
    providers:[
        CharlesService
    ]
})
export class CharlesModule{}