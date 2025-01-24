import styles from "./MovieModal.module.scss";
import InputSelect from "../../inputs/inputSelect/InputSelect";
import InputCheckbox from "../../inputs/inputCheckbox/InputCheckbox";
import Modal from "../../modal/Modal";

const MovieModal = ({
	isOpen,
	setIsOpen,
	content,
	setContent,
	profiles,
}: any) => {
	const onSave = async () => {
		await fetch(`/api/movies/${content.id}`, {
			method: "PUT",
			headers: {
				"Content-Type": "application/json",
				Authorization: `Bearer ${localStorage.getItem("token")}`,
			},
			body: JSON.stringify(content),
		});
		setIsOpen(false);
	};

	const onClose = () => {
		setIsOpen(false);
	};
	if (!isOpen) return null;
	return (
		<Modal
			isOpen={isOpen}
			setIsOpen={setIsOpen}
			title={"Movie Options"}
			onClose={onClose}
			onSave={onSave}
		>
			<div className={styles.content}>
				<div className={styles.inputContainer}>
					<label className={styles.label}>Monitored </label>
					<InputCheckbox
						type="checkbox"
						checked={content?.monitored}
						onChange={(e: any) =>
							setContent({
								...content,
								monitored: e.target.checked,
							})
						}
					/>
				</div>
				<div className={styles.inputContainer}>
					<label className={styles.label}>Profile </label>
					<InputSelect
						selected={content.profileId}
						onChange={(e: any) => {
							setContent({ ...content, profileId: parseInt(e.target.value) });
						}}
					>
						{Object.values(profiles)?.map((profile: any, index: number) => (
							<option value={profile.id} key={index}>
								{profile.name}
							</option>
						))}
					</InputSelect>
				</div>
			</div>
		</Modal>
	);
};
export default MovieModal;
