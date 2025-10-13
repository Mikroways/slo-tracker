import React, { useEffect, useState } from 'react';
import { Button, Form, Input, Checkbox, Row, Col, Popconfirm, TimePicker } from 'antd';
import moment from 'moment';
import dayjs from 'dayjs';
import SLOService from '../../../core/services/service.slo';
import { ISLO } from '../../../core/interfaces/ISLO';
import openNotification from '../../../core/helpers/notification';

interface IProps {
  refreshSLOs: () => void;
  closeDrawer: () => void;
  activeSLO: ISLO | null;
  open: boolean;
}

const DAYS = ["sunday","monday","tuesday","wednesday","thursday","friday","saturday"];

const UpdateSLO: React.FC<IProps> = (props) => {
  const _sloService = new SLOService();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  // Fetch SLO details when activeSLO changes
  useEffect(() => {
    async function fetchSLO() {
      if (props.open && props.activeSLO) {
        setLoading(true);
        try {
          const response = await _sloService.getWorkingSchedule(props.activeSLO.id);
          const sloArray = response.data.data;
          const slo = Array.isArray(sloArray) ? sloArray[0] : sloArray;

          // Prepare initial values for working days
          const initialDays: any = {};
          if (slo.working_days) {
            slo.working_days.forEach((wd: any) => {
              const day = DAYS[wd.weekday];
              if (!day) return;

              initialDays[`${day}_enabled`] = true;
              initialDays[`${day}_open_hour`] = wd.open_hour ? moment(wd.open_hour, "HH:mm") : null;
              initialDays[`${day}_close_hour`] = wd.close_hour ? moment(wd.close_hour, "HH:mm") : null;
            });
          }
          DAYS.forEach(day => {
            if (!initialDays.hasOwnProperty(`${day}_enabled`)) {
              initialDays[`${day}_enabled`] = false;
              initialDays[`${day}_open_hour`] = null;
              initialDays[`${day}_close_hour`] = null;
            }
          });

          form.resetFields();

          form.setFieldsValue({
            slo_name: slo.slo_name,
            holidays_enabled: slo.holidays_enabled,
            target_slo: slo.target_slo,
            reset_slo: slo.holidays_enabled,
            ...initialDays,
          });
        } catch (err) {
          openNotification('error', 'Failed to fetch SLO details.');
        } finally {
          setLoading(false);
        }
      } else if (!props.open){
        form.resetFields();
      }
    }
    fetchSLO();
  }, [props.activeSLO, props.open, form]);

  const onSubmit = async (values: any) => {
    if (!props.activeSLO) return;

    const slo_name = values['slo_name'];
    const target_slo = parseFloat(values['target_slo']);
    const reset_slo = values['reset_slo'] || false;
    const holidays_enabled = values['holidays_enabled'];

    // Validate enabled days first
    for (const day of DAYS) {
      if (values[`${day}_enabled`]) {
        const open_hour = values[`${day}_open_hour`];
        const close_hour = values[`${day}_close_hour`];
        if (!open_hour || !close_hour) {
          openNotification('error', `Please provide opening and closing hours for ${day}.`);
          return;
        }
        if (open_hour.format("HH:mm") >= close_hour.format("HH:mm")) {
          openNotification('error', `Closing Hour should be after Opening Hour for ${day}.`);
          return;
        }
      }
    }

    // Only include enabled days
    const working_days = DAYS.filter(day => values[`${day}_enabled`]).map(day => ({
      'weekday': DAYS.indexOf(day),
      open_hour: values[`${day}_open_hour`].format("HH:mm"),
      close_hour: values[`${day}_close_hour`].format("HH:mm"),
    }));

    try {
      await _sloService.update(
        props.activeSLO.id,
        {
          slo_name,
          target_slo,
          working_days,
          'holidays_enabled': holidays_enabled,
        },
        reset_slo
      );
      props.refreshSLOs();
      openNotification('success', 'Successfully updated SLO');
      props.closeDrawer();
      form.resetFields();
    } catch (err) {
      openNotification('error', 'Error while updating SLO. Please try again.');
    }
  };

  const onConfirmDelete = async () => {
    if (!props.activeSLO) return;

    try {
      await _sloService.delete(props.activeSLO.id);

      props.refreshSLOs();
      openNotification('success', 'Successfully Deleted SLO');
      props.closeDrawer();
      form.resetFields();
    } catch (err) {
      openNotification('error', 'Error while deleting SLO. Please try again.');
    }
  };

  return (
    <Form
      layout="vertical"
      onFinish={onSubmit}
      form={form}
    >
      <Form.Item
        label="SLO Name"
        name="slo_name"
        rules={[{ required: true, message: 'Please give your SLO a name!' }]}
      >
        <Input placeholder="Eg: Checkout Flow" />
      </Form.Item>

      <Form.Item
        label="Target SLO in %"
        name="target_slo"
        rules={[{ required: true, message: 'Please provide a target SLO' }]}
      >
        <Input placeholder="Eg: 99.999" />
      </Form.Item>

      <Form.Item name="reset_slo" valuePropName="checked">
        <Checkbox>Reset complete Error-budget</Checkbox>
      </Form.Item>

      <div style={{ marginBottom: 16, fontWeight: 500 }}>Days & Hours</div>
      {DAYS.map(day => (
        <Row key={day} gutter={8} align="middle" style={{ marginBottom: 8 }}>
          <Col style={{ minWidth: 120 }}>
            <Form.Item name={`${day}_enabled`} valuePropName="checked" noStyle>
              <Checkbox
              onChange={e => {
                if (e.target.checked) {
                  form.setFieldsValue({
                    [`${day}_open_hour`]: dayjs().hour(0).minute(0),
                    [`${day}_close_hour`]: dayjs().hour(23).minute(59),
                  });
                } else {
                  form.setFieldsValue({
                    [`${day}_open_hour`]: null,
                    [`${day}_close_hour`]: null,
                  });
                }
              }}
              >{day.charAt(0).toUpperCase() + day.slice(1)}</Checkbox>
            </Form.Item>
          </Col>
          <Col>
            <Form.Item shouldUpdate={(prev, curr) => prev[`${day}_enabled`] !== curr[`${day}_enabled`]}>
              {({ getFieldValue }) => (
                <Form.Item
                  name={`${day}_open_hour`}
                  noStyle
                  rules={[
                    {
                      validator(_, value) {
                        if (!getFieldValue(`${day}_enabled`) || value) return Promise.resolve();
                        return Promise.reject(new Error('Please provide opening hour'));
                      }
                    }
                  ]}
                >
                  <TimePicker
                    format="HH:mm"
                    minuteStep={15}
                    placeholder="Open"
                    disabled={!getFieldValue(`${day}_enabled`)}
                  />
                </Form.Item>
              )}
            </Form.Item>
          </Col>
          <Col>
            <Form.Item shouldUpdate={(prev, curr) => prev[`${day}_enabled`] !== curr[`${day}_enabled`]}>
              {({ getFieldValue }) => (
                <Form.Item
                  name={`${day}_close_hour`}
                  noStyle
                  rules={[
                    {
                      validator(_, value) {
                        if (!getFieldValue(`${day}_enabled`) || value) return Promise.resolve();
                        return Promise.reject(new Error('Please provide closing hour'));
                      }
                    }
                  ]}
                >
                  <TimePicker
                    format="HH:mm"
                    minuteStep={15}
                    placeholder="Close"
                    disabled={!getFieldValue(`${day}_enabled`)}
                  />
                </Form.Item>
              )}
            </Form.Item>
          </Col>
        </Row>
      ))}

      <Form.Item name="holidays_enabled" valuePropName="checked">
        <Checkbox>
          Enable Holidays
        </Checkbox>
      </Form.Item>

      <Form.Item>
        <Row style={{ justifyContent: 'space-between' }}>
          <Button type="primary" htmlType="submit" loading={loading}>
            Update SLO
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this SLO?"
            okText="Yes. Delete SLO"
            placement="topRight"
            onConfirm={onConfirmDelete}
          >
            <Button danger>Delete SLO</Button>
          </Popconfirm>
        </Row>
      </Form.Item>
    </Form>
  );
};

export default UpdateSLO;