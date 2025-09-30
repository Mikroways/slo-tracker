import { DatePicker, Table, Button } from "antd";
import React, { useState } from "react";
import { Link } from "react-router-dom";
import moment, { Moment } from 'moment';

import useGetOverview from "../../../core/hooks/userGetOverview";
import { ISLO } from "../../../core/interfaces/ISLO";


interface IProps {
	setActiveSLO: (activeSLO: ISLO) => void;
}

const tableColumns = (props) => [
	{
		title: "SLO",
		key: "slo_name",
		render: (o: ISLO) => (
			<Link
				to="/"
				onClick={() => {
					props.setActiveSLO(o);
				}}
				style={{ color: "#0f0f0f" }}
			>
				{o.slo_name}
			</Link>
		),
	},
	{
		title: "Target",
		dataIndex: "target_slo",
		key: "target_slo",
	},
	{
		title: "Your SLO",
		render: (o: ISLO) => (o.current_slo === 0 ? 100 : o.current_slo),
	},
	{
		title: "Remaining Error Budget (mins)",
		dataIndex: "remaining_err_budget",
		key: "remaining_err_budget",
	},
	{
		title: "Incidents Reported",
		dataIndex: "num_incidents",
		key: "num_incidents"
	},
	{
		title: "Incidents Reported (False)",
		dataIndex: "num_incidents_false_positive",
		key: "num_incidents_false_positive"
	},
	{
		title: "View Incidents",
		key: "View Incidents",
		render: (o: ISLO) => (
			<Link
				to="/"
				onClick={() => {
					props.setActiveSLO(o);
				}}
			>
				View Incidents
			</Link>
		),
	},
];

const DownloadCSVButton = ({ data, yearMonth }: { data: any[]; yearMonth: string }) => {
  const downloadCSV = () => {
	const headers = [
	  "SLO Name",
	  "Target SLO",
	  "Current SLO",
	  "Remaining Error Budget",
	  "Incidents",
	  "False Positives",
	];

    const rows = data.map((row) => [
      row.slo_name,
      row.target_slo,
      row.current_slo,
      row.remaining_err_budget,
      row.num_incidents,
      row.num_incidents_false_positive,
    ]);

    const csvContent =
      [headers, ...rows].map((e) => e.join(",")).join("\n");

    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);

    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", `overview-${yearMonth}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <Button type="primary" onClick={downloadCSV}>
      Download CSV
    </Button>
  );
};

const SLOTable: React.FC<IProps> = ({ ...props }) => {
	const [selectedMonth, setSelectedMonth] = useState<Moment>(moment());
	  const handleMonthChange = (value: Moment | null, dateString: string) => {
		setSelectedMonth(value ?? moment());
	  };

	const { overview } = useGetOverview(
		selectedMonth ? selectedMonth.format('YYYY-MM') : undefined
	);

	const tableData = overview.map((i) => ({
		key: i.id,
		...i,
	}));
	return (
		<>
			<div style={{ fontSize: 20, fontWeight: 300, margin: "20px 0", gap: "0.7em", display: "flex", alignItems: "center" }}>
				<Link to="/" style={{ marginLeft: "1em" }}>
					SLO Tracker
				</Link>
				{" / Overview"}
				<DatePicker
				picker="month"
				format="MMMM YYYY"
				value={selectedMonth}
				onChange={handleMonthChange}
				/>
				<DownloadCSVButton
					data={tableData}
					yearMonth={selectedMonth.format("YYYY-MM")}
				/>
			</div>
			<Table
				dataSource={tableData}
				columns={tableColumns(props)}
				pagination={{ pageSize: 15 }}
				size="middle"
				style={{ width: "auto", margin: "0 2em" }}
			/>
		</>
	);
};

export default SLOTable;
