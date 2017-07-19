import { BasePage } from './../base.page';
import { CharlesService } from './../../app/charles/charles.service';
import { CharlesModel } from './../../app/charles/charles.model';
import { Component } from '@angular/core';
import { NavController, AlertController, LoadingController } from 'ionic-angular';

@Component({
  selector: 'page-contact',
  templateUrl: 'charles.html'
})
export class CharlesPage extends BasePage {
  charlesModel = new CharlesModel();

  constructor(
    alertCtrl: AlertController,
    navCtrl: NavController,
    loadingCtrl: LoadingController,
    private charlesService: CharlesService
  ) {
    super(alertCtrl, navCtrl, loadingCtrl);
  }

  getKey() {
    this.execute(function (parent) {
      parent.charlesService.getKey(parent.charlesModel)
        .then(response => {
          let model = (<any>response).json() as CharlesModel;
          parent.charlesModel.key = model.key;
        })
        .then(() => { parent.dismissLoader() });
    })
  }
}
