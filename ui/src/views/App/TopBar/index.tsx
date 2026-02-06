import React from 'react';
import { Button } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import Overview from '../Overview/Overview';

interface IProps {
	onAddSLO: () => void;
}

const TopBar: React.FC<IProps> = (props) => {
	return (
		<div style={{ display: 'flex', justifyContent: 'flex-start', alignItems: 'center' }}>
			<Overview />
			<Button type="dashed" icon={<PlusOutlined/>} onClick={props.onAddSLO}>
				Create SLO
			</Button>
		</div>
	);
};

export default TopBar;
