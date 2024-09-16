import styles from "./Profile.module.scss";
import AddIcon from "../svgs/add.svg?react";
import Codec from "../codec/Codec";

const Profile = ({ name, type = "", codecs, onClick, profile }: any) => {
	return (
		<div className={styles.profile} key={name} onClick={() => onClick(profile)}>
			{type === "add" ? (
				<div className={styles.add}>
					<div className={styles.box}>
						<AddIcon className={styles.svg} />
					</div>
				</div>
			) : (
				<div className={styles.normal}>
					<div className={styles.header}>{name}</div>
					<div className={styles.codecs}>
						{codecs.map((codec: any) => (
							<Codec key={"codec" + codec.codecId} codec={codec.codecId} />
						))}
					</div>
				</div>
			)}
		</div>
	);
};
export default Profile;
