import { AppStorePage } from './../appstore/appstore';
import { CharlesPage } from './../charles/charles';
import { EditPlusPage } from './../editplus/editplus';
import { Component } from '@angular/core';

@Component({
  templateUrl: 'tabs.html'
})
export class TabsPage {

  tab1Root = EditPlusPage;
  tab2Root = CharlesPage;
  tab3Root = AppStorePage;

  constructor() {

  }
}
