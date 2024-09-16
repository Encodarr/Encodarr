import styles from "./ProfileSummary.module.scss";
import InputContainer from "../inputs/inputContainer/InputContainer";
import InputCheckbox from "../inputs/inputCheckbox/InputCheckbox";

const ProfileSummary = ({ content, setContent, containers, codecs }: any) => {
	console.log(content)
	return (
		<div className={styles.section}>
			<div className={styles.left}>
				<InputContainer
					type="text"
					label="Profile Name"
					selected={content.name}
					onChange={(e: any) =>
						setContent({ ...content, name: e.target.value })
					}
				/>
				<InputContainer
					type="select"
					label="Output Container"
					selected={content.container}
					onChange={(e: any) => {
						setContent({
							...content,
							container: e.target.value,
							extension: containers[e.target.value].extensions[0],
						});
					}}
				>
					{codecs[content?.codec]?.containers?.map(
						(container: any, index: number) => (
							<option value={container} key={index}>
								{container}
							</option>
						)
					)}
				</InputContainer>
				<InputContainer
					type="select"
					label="Output Extension"
					selected={content.extension}
					onChange={(e: any) =>
						setContent({ ...content, extension: e.target.value })
					}
				>
					{containers[content?.container]?.extensions?.map(
						(extension: any, index: number) => (
							<option value={extension} key={index}>
								.{extension}
							</option>
						)
					)}
				</InputContainer>
				<InputContainer
					type="checkbox"
					label=""
					checked={content.passThruCommonMetadata
					}
					helpText="Passthru Common Metadata"
					onChange={(e: any) =>
						setContent({
							...content,
							passThruCommonMetadata
								: e.target.checked,
						})
					}
				/>
				<div className={styles.description}></div>
			</div>
			<div className={styles.right}>
				<label>Targets</label>
				<div className={styles.targets}>
					{Object.entries(codecs)?.map(([key]: any) => (
						<div
							key={key}
							className={styles.target}
							style={!content?.codecs?.some(codec => codec.codecId === key) ? { opacity: "50%" } : {}} onClick={() => {
								const selectedValues = [...content.codecs];
								if (selectedValues.includes(key)) {
									const index = selectedValues.indexOf(key);
									if (index > -1) {
										selectedValues.splice(index, 1);
									}
								} else {
									selectedValues.push({ codecId: key, profileId: content.id });
								}
								setContent({ ...content, codecs: selectedValues });
							}}
						>
							<InputCheckbox
								type="checkbox"
								checked={content?.codecs?.some(codec => codec.codecId === key)}
								onChange={() => {
									const selectedValues = [...content.codecs];
									const index = selectedValues.findIndex(codec => codec.codecId === key);
									if (index > -1) {
										selectedValues.splice(index, 1);
									} else {
										selectedValues.push({ codecId: key });
									}
									setContent({ ...content, codecs: selectedValues });
								}}
							/>
							<span className={styles.key}>{key}</span>
						</div>
					))}
				</div>
			</div>
		</div>
	);
};
export default ProfileSummary;
