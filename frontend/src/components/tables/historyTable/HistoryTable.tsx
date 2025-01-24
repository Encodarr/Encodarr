import styles from "./HistoryTable.module.scss";
import ResetWrench from "../../svgs/reset_wrench.svg?react";
import Codec from "../../codec/Codec";
import { formatDate, formatSize } from "../../../utils/format";
import Table from "../../table/Table";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

const HistoryTable = ({ historyArray, settings }: any) => {
	const recordsPerPage = settings?.historyPageSize || 0;
	const [currentPage, setCurrentPage] = useState(1);
	const indexOfLastRecord = currentPage * recordsPerPage;
	const indexOfFirstRecord = indexOfLastRecord - recordsPerPage;
	const currentRecords = historyArray.slice(
		indexOfFirstRecord,
		indexOfLastRecord
	);
	const totalPages = Math.ceil(historyArray.length / recordsPerPage);
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
			totalRecords={historyArray.length}
			lastPage={lastPage}
			nextPage={nextPage}
		>
			<thead>
				<tr className={styles.headRow}>
					<th></th>
					<th>Title</th>
					<th>Type</th>
					<th>Episode</th>
					<th>Codec</th>
					<th>Original Codec</th>
					<th>Space Impact</th>
					<th>Date</th>
				</tr>
			</thead>
			<tbody>
				{currentRecords.reverse().map((entry: any, index: number) => (
					<tr className={styles.row} key={index}>
						<td>
							<ResetWrench className={styles.svg} />
						</td>
						<td>
							<Link
								to={
									(entry?.type == "episode" ? "/series/" : "/movies/") +
									entry?.name
								}
								className={styles.name}
							>
								{entry?.name}
							</Link>
						</td>
						<td>{entry?.type == "movie" ? "Movie" : "Series"}</td>
						<td>
							{entry?.type == "episode" && (
								<>
									{entry?.seasonNumber}x{entry?.episodeNumber}
								</>
							)}
						</td>
						<td className={styles.codecRow}>
							<Codec codec={entry?.prevCodec} />
						</td>
						<td className={styles.codecRow}>
							<Codec codec={entry?.newCodec} />
						</td>
						<td>{formatSize(entry?.prevSize - entry?.newSize)}</td>
						<td>{formatDate(entry?.date)}</td>
					</tr>
				))}
			</tbody>
		</Table>
	);
};
export default HistoryTable;
