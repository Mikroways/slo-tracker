import React from 'react';
import { Button, Form, Input, TimePicker, Row, Col, Checkbox } from 'antd';
import dayjs from 'dayjs';
import SLOService from '../../../core/services/service.slo';
import openNotification from '../../../core/helpers/notification';

interface IProps {
  refreshSLOs: () => void;
  closeDrawer: () => void;
}

const DAYS = ["sunday","monday","tuesday","wednesday","thursday","friday","saturday"];

const CreateSLO: React.FC<IProps> = (props) => {
  const _sloService = new SLOService();
  const [form] = Form.useForm();

  const onSubmit = async (values: any) => {
    const slo_name = values['slo_name'];
    const target_slo = parseFloat(values['target_slo']);
    const holidays_enabled = values['holidays_enabled']

    if (target_slo < 1 || target_slo > 100) {
      openNotification('error', 'Target SLO should be between 1 to 100.');
      return;
    }

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
      await _sloService.create({
        slo_name,
        target_slo,
        working_days,
        holidays_enabled,
      });
      props.refreshSLOs();
      openNotification('success', 'Successfully created SLO');
      props.closeDrawer();
      form.resetFields();
    } catch (err) {
      openNotification('error', 'Error while creating SLO. Please try again.');
    }
  };

  return (
    <Form layout="vertical" onFinish={onSubmit} form={form}>
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

      <div style={{ marginBottom: 16, fontWeight: 500 }}>Days & Hours</div>
      {DAYS.map(day => (
        <Row key={day} gutter={8} align="middle"  justify="start" style={{ marginBottom: 8 }}>
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
              >  
                {day.charAt(0).toUpperCase() + day.slice(1)}
              </Checkbox>
            </Form.Item>
          </Col>
          <Col>
            <Form.Item shouldUpdate={(prev, curr) => prev[`${day}_enabled`] !== curr[`${day}_enabled`]}>
              {({ getFieldValue }) => (
                <Form.Item
                  name={`${day}_open_hour`}
                  rules={[
                    {
                      validator(_, value) {
                        if (!getFieldValue(`${day}_enabled`) || value) return Promise.resolve();
                        return Promise.reject(new Error('Please provide opening hour'));
                      }
                    }
                  ]}
                  noStyle
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
                  rules={[
                    {
                      validator(_, value) {
                        if (!getFieldValue(`${day}_enabled`) || value) return Promise.resolve();
                        return Promise.reject(new Error('Please provide closing hour'));
                      }
                    }
                  ]}
                  noStyle
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

      <Form.Item name="holidays_enabled" valuePropName="checked" noStyle>
        <Checkbox>
          Enable Holidays
        </Checkbox>
      </Form.Item>

      <Form.Item>
        <Button style={{ float: 'right' }} type="primary" htmlType="submit">
          Create SLO
        </Button>
      </Form.Item>
    </Form>
  );
};

export default CreateSLO;
