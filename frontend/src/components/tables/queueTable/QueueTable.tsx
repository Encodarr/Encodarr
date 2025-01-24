import styles from "./QueueTable.module.scss";
import QueueIcon from "../../svgs/queue.svg?react";
import Codec from "../../codec/Codec";
import { formatSize } from "../../../utils/format";
import Table from "../../table/Table";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
const QueueTable = ({ queueArray, profiles, series, settings }: any) => {

	const [currentPage, setCurrentPage] = useState(1);
	series = series ? () => {
		return series?.reduce((map, item) => {
			map[item.id] = item;
			return map;
		}, {});
	} : {};

	const recordsPerPage = settings?.queuePageSize || 0;
	const totalPages = Math.ceil(queueArray?.length / recordsPerPage);
	const indexOfLastRecord = currentPage * recordsPerPage;
	const indexOfFirstRecord = indexOfLastRecord - recordsPerPage;
	const currentRecords = queueArray?.slice(
		indexOfFirstRecord,
		indexOfLastRecord
	);
	const firstPage = () => {
		setCurrentPage(1);
	};
	const lastPage = () => {
		setCurrentPage(totalPages);
	};
	const nextPage = () => {
		if (currentPage < totalPages) {
			setCurrentPage(currentPage + 1);
		}
	};
	const prevPage = () => {
		if (currentPage > 1) {
			setCurrentPage(currentPage - 1);
		}
	};
	useEffect(() => {
		setCurrentPage(1);
	}, [recordsPerPage]);
	return (
		<Table
			showPagination={true}
			firstPage={firstPage}
			currentPage={currentPage}
			prevPage={prevPage}
			totalPages={totalPages}
			totalRecords={queueArray?.length}
			lastPage={lastPage}
			nextPage={nextPage}
		>
			<thead>
				<tr>
					<th></th>
					<th>Title</th>
					<th>Type</th>
					<th>Episode</th>
					<th>Profile</th>
					<th>Codec</th>
					<th>Future Codec</th>
					<th>Size</th>
				</tr>
			</thead>
			<tbody>
				{currentRecords?.map((q: any) => (
					<tr>
						<td>
							<QueueIcon className={styles.svg} />
						</td>
						<td>
							<Link to={
								q?.seriesId ? "/series/" + q?.seriesId : "/movies/" + q?.id
							} className={styles.name}>

								{q?.seriesId ? q?.seriesId : q?.id}
							</Link>
						</td>
						<td>{q?.type == "episode" ? "Series" : "Movie"}</td>
						<td>
							{q?.type == "episode" && (
								<>
									{q.seasonNumber}x{q.episodeNumber}
								</>
							)}
						</td>
						<td>
							{profiles && series && (
								profiles.find((profile: any) => profile.id == q?.profileId).name ?? "Profile not found"
							)}
						</td>
						<td className={styles.codecRow}>
							<Codec codec={q.codec} />
						</td>
						<td className={styles.codecRow}>
							<Codec
								codec={profiles && series && (
									profiles.find((profile: any) => profile.id == q?.profileId).codec ?? "Profile not found"
								)}
							/>
						</td>
						<td>{formatSize(q?.size)}</td>
					</tr>
				))}
			</tbody>
		</Table>
	);
};
export default QueueTable;
