import { useContext } from "react";
import styles from "./Footer.module.scss";
import { WebSocketContext } from "../../contexts/webSocketContext";
import { formatSize } from "../../utils/format";
const Footer = () => {
	const wsContext = useContext(WebSocketContext);
	const systemArray: any = wsContext?.data?.system;
	
	const systemDictionary: { [key: string]: any } = systemArray?.length === 0 
  ? {} 
  : systemArray.reduce((acc, obj) => {
      acc[obj.id] = obj;
      return acc;
    }, {});
	return (
		<div className={styles.footer}>
			<div className={styles.keyContainer}>
				<div className={styles.key}>
					<div
						className={styles.bubble}
						style={{ backgroundColor: "#5d9cec" }}
					></div>
					<div>Continuing (All episodes downloaded)</div>
				</div>
				<div className={styles.key}>
					<div
						className={styles.bubble}
						style={{ backgroundColor: "#27c24c" }}
					></div>
					<div>Ended (All episodes downloaded)</div>
				</div>
				<div className={styles.key}>
					<div
						className={styles.bubble}
						style={{ backgroundColor: "#f05050" }}
					></div>
					<div>Missing Episodes (Series monitored)</div>
				</div>
				<div className={styles.key}>
					<div
						className={styles.bubble}
						style={{ backgroundColor: "#ffa500" }}
					></div>
					<div>Missing Episodes (Series not monitored)</div>
				</div>
			</div>
			<div className={styles.dataContainer}>
				<div className={styles.section}>
					<div className={styles.data}>Series {systemDictionary?.series_count?.value}</div>
					<div className={styles.data}>Ended {systemDictionary?.ended_count?.value}</div>
					<div className={styles.data}>
						Continuing {systemDictionary?.continuing_count?.value}
					</div>
				</div>
				<div className={styles.section}>
					<div className={styles.data}>Monitored {systemDictionary?.monitored_count?.value}</div>

					<div className={styles.data}>
						Unmonitored {systemDictionary?.unmonitored_count?.value}
					</div>
				</div>
				<div className={styles.section}>
					<div className={styles.data}>Episodes {systemDictionary?.episode_count?.value}</div>
					<div className={styles.data}>Files {systemDictionary?.files_count?.value}</div>
				</div>
				<div className={styles.section}>
					<div className={styles.data}>
						{"Size on Disk   "}
						{formatSize(systemDictionary?.size_on_disk?.value)}
					</div>
					<div className={styles.data}>
						{"Space Saved   "}
						{formatSize(systemDictionary?.space_saved?.value)}
					</div>
				</div>
			</div>
		</div>
	);
};
export default Footer;
