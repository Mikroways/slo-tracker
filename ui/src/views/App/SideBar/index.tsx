import React, { useMemo, useState } from 'react';
import { Menu, Input, Tooltip } from 'antd';

import { ISLO } from '../../../core/interfaces/ISLO';

import './sidebar.css';

interface IProps {
  SLOs: ISLO[];
  activeSLOId?: number;
  setActiveSlo: (slo: ISLO) => void;
}

const SideBar: React.FC<IProps> = ({ SLOs, activeSLOId, ...props }) => {
  const [filter, setFilter] = useState('');

  const onSLOSelect = (value: any) => {
    const SLOId = parseInt(value.key);
    const selectedSLO = SLOs.filter((s) => s.id === SLOId);

    if (selectedSLO.length) props.setActiveSlo(selectedSLO[0]);
  };

  const filteredSLOs = useMemo(() => {
    const q = filter.trim().toLowerCase();
    if (!q) return SLOs;
    return SLOs.filter((s) => s.slo_name.toLowerCase().includes(q));
  }, [SLOs, filter]);

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'flex-start',
        height: '100%',
      }}
    >
      <div>
        <Input.Search
          placeholder="Filter SLOs"
          allowClear
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          style={{ marginBottom: 8 }}
        />

        <Menu defaultSelectedKeys={[String(activeSLOId)]} onClick={onSLOSelect}>
          {filteredSLOs.map((slo) => (
            <Menu.Item key={String(slo.id)}>
              <Tooltip
                title={
                  <div style={{ textAlign: 'left' }}>
                    <div><strong>{slo.slo_name}</strong></div>
                  </div>
                }
                placement="right"
              >
                <span>{slo.slo_name}</span>
              </Tooltip>
            </Menu.Item>
          ))}
        </Menu>
      </div>
    </div>
  );
};

export default SideBar;
