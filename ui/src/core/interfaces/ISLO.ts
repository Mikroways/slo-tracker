export interface ISLO {
  id: number;
  slo_name: string;
  target_slo: number;
  current_slo: number;
  updated_at: Date;
  remaining_err_budget: number;
  working_days?: any;
  holidays_enabled: boolean;
  num_incidents?: number;
  num_incidents_false_positive?: number;
}
