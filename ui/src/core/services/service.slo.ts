import axios from 'axios';
import { API } from '.';
import { IResponse } from '../interfaces/IResponse';
import { ISLO } from '../interfaces/ISLO';

class SLOService {
  private _api = `${API}/api/v1/slo`;

  all = () => axios.get<IResponse<ISLO[]>>(this._api);

  get = (sloId: number, params?: {yearMonth?: string}) => axios.get<IResponse<ISLO>>(`${this._api}/${sloId}`, { params });

  getWorkingSchedule = (sloId: number) => axios.get<IResponse<ISLO[]>>(`${this._api}/${sloId}/workingschedule`);

  create = (data: Pick<ISLO, 'slo_name' | 'target_slo' | 'working_days' | 'holidays_enabled'>) =>
    axios.post<IResponse<ISLO>>(this._api, data);

  update = (
    sloId: number,
    data: Pick<ISLO, 'slo_name' | 'target_slo' | 'working_days' | 'holidays_enabled'>,
    isReset: boolean
  ) =>
    axios.patch<IResponse<ISLO>>(
      `${this._api}/${sloId}?isReset=${isReset}`,
      data
    );

  delete = (sloId: number) => axios.delete(`${this._api}/${sloId}`);
}

export default SLOService;
