import { Switch, Table, Typography } from 'antd';
import React from 'react';
import openNotification from '../../../../core/helpers/notification';
import { IIncident } from '../../../../core/interfaces/IIncident';
import IncidentService from '../../../../core/services/service.incident';

const { Paragraph } = Typography

interface IProps {
  SLIs: IIncident[];
  refreshSLOAndSLIs: (slo?: boolean, slis?: boolean) => void;
}

const tableColumns = (
  updateObservations: (sli: IIncident, newText: string) => void,
  onMarkPositive: (state: boolean, sli: IIncident) => void
) => [
  {
    title: 'SLI',
    dataIndex: 'sli_name',
    key: 'sli_name',
  },
  {
    title: 'Status',
    dataIndex: 'state',
    key: 'state',
  },
  {
    title: 'Alert Source',
    dataIndex: 'alertsource',
    key: 'alertsource',
  },
  {
    title: 'Created on',
    key: 'created_at',
    render: (e: IIncident) => {
      const date = new Date(e.created_at);
      return <p>{date.toLocaleString()}</p>;
    },
  },
  {
    title: 'Error budget spent(min)',
    dataIndex: 'err_budget_spent',
    key: 'err_budget_spent',
  },
  {
    title: 'Real Err budget spent(min)',
    dataIndex: 'real_err_budget_spent',
    key: 'real_err_budget_spent',
  },
  {
    title: 'Observations',
    dataIndex: 'observations',
    key: 'observations',
    render: (value: string, sli: IIncident) => (
        <Paragraph
          editable={{
            onChange: (newText) => updateObservations(sli, newText),
            tooltip: "Click to edit observations",
          }}
        >
          {value}
        </Paragraph>
    ),
  },
  {
    title: 'Mark false positive',
    key: 'action',
    render: (sli: IIncident) => (
      <Switch
        onChange={(state: boolean) => onMarkPositive(state, sli)}
        defaultChecked={sli.mark_false_positive}
      />
    ),
  },
];

const SLITable: React.FC<IProps> = ({ SLIs, ...props }) => {

  const [data, setData] = React.useState<IIncident[]>(SLIs);

  React.useEffect(() => {
    setData(SLIs);
  }, [SLIs]);

  const tableData = data.map((i) => ({
    key: i.id,
    ...i,
  }));

  const updateObservations = async (sli: IIncident, newText: string) => {

    const _incidentService = new IncidentService(sli.id);

    if (newText.length >= 500) {
      openNotification('error', 'Observations must be less than 500 characters');
      return;
    }

    try {
      await _incidentService.update(sli.id, {
        observations: newText,
        mark_false_positive: sli.mark_false_positive,
        state: sli.state,
        err_budget_spent: sli.err_budget_spent,
        real_err_budget_spent: sli.real_err_budget_spent,
      });
      props.refreshSLOAndSLIs();
      openNotification('success', 'Incident updated successfully');
    } catch (err) {
      console.log(err);
      openNotification('error', 'Error while updating incident');
    }

  };

  const onMarkPositive = async (state: boolean, sli: IIncident) => {
    const _incidentService = new IncidentService(sli.slo_id);

    try {
      await _incidentService.update(sli.id, {
        mark_false_positive: state,
        state: sli.state,
        err_budget_spent: sli.err_budget_spent,
        real_err_budget_spent: sli.real_err_budget_spent,
        observations: sli.observations,
      });
      props.refreshSLOAndSLIs();
      openNotification('success', 'Incident updated successfully');
    } catch (err) {
      console.log(err);
      openNotification('error', 'Error while updating incident');
    }
  };

  return (
    <Table
      dataSource={tableData}
      columns={tableColumns(updateObservations, onMarkPositive)}
      pagination={{ pageSize: 5 }}
    />
  );
};

export default SLITable;
