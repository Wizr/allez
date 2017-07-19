import { CharlesModel } from './charles.model';
import { Http } from '@angular/http';
import { GlobalVariables } from './../global.variables';
import { BaseService } from './../base.service';
import { Injectable } from '@angular/core';

@Injectable()
export class CharlesService extends BaseService {
    private charlesUrl = GlobalVariables.API_URL + 'api/keygen';

    constructor(
        public http: Http
    ) {
        super(http);
    }

    getKey(model: CharlesModel) {
        let url = `${this.charlesUrl}/charles`;
        return this.post(url, false, JSON.stringify(model), this);
    }
}
