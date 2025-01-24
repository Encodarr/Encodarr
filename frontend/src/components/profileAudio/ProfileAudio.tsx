import InputCheckbox from "../inputs/inputCheckbox/InputCheckbox";
import InputContainer from "../inputs/inputContainer/InputContainer";
import styles from "./ProfileAudio.module.scss";
const ProfileAudio = ({ content, setContent }) => {
	const languages = {
		all: "All",
		eng: "English",
		spa: "Spanish",
		fra: "French",
		deu: "German",
		ita: "Italian",
		jpn: "Japanese",
		kor: "Korean",
		rus: "Russian",
		zho: "Chinese",
	};
	return (
		<>
			<div className={styles.section}>
				<div className={styles.left}>
					<InputContainer
						label="Include Untagged Audio Tracks"
						type="checkbox"
						checked={content?.mapUntaggedAudioTracks}
						onChange={(e) => {
							setContent({
								...content,
								mapUntaggedAudioTracks: e.target.checked,
							});
						}}
					/>
				</div>
				<div className={styles.right}>
					<label>Wanted Audio Tracks</label>
					<div className={styles.targets}>
						{Object.entries(languages)?.map(([key, value]: any) => (
							<div
								key={key}
								className={styles.target}
								style={
									!content?.profileAudioLanguages?.some(lang => lang.language === key)
										? { opacity: "50%" }
										: {}
								}
								onClick={() => {
									const selectedValues = [...content.profileAudioLanguages];
									const index = selectedValues.findIndex(lang => lang.language === key);
									if (index > -1) {
										selectedValues.splice(index, 1);
									} else {
										selectedValues.push({ profileId: content.id, language: key });
									}
									setContent({ ...content, profileAudioLanguages: selectedValues });
								}}
							>
								<InputCheckbox
									type="checkbox"
									checked={content?.profileAudioLanguages?.some(lang => lang.language === key)}
									onChange={() => {
										const selectedValues = [...content.profileAudioLanguages];
										const index = selectedValues.findIndex(lang => lang.language === key);
										if (index > -1) {
											selectedValues.splice(index, 1);
										} else {
											selectedValues.push({ profileId: content.id, language: key });
										}
										setContent({ ...content, profileAudioLanguages: selectedValues });
									}}
								/>
								<span className={styles.key}>{value}</span>
							</div>
						))}
					</div>
				</div>
			</div>
		</>
	);
};
export default ProfileAudio;
