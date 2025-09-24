export interface ISLO {
  id: number;
  slo_name: string;
  open_hour?: Date,
  close_hour?: Date,
  target_slo: number;
  current_slo: number;
  updated_at: Date;
  remaining_err_budget: number;
  working_days?: any;
}
