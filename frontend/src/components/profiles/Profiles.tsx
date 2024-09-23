import styles from "./Profiles.module.scss";
import Profile from "../profile/Profile";
import { useContext, useState } from "react";
import ProfileModal from "../modals/profileModal/ProfileModal";
import { WebSocketContext } from "../../contexts/webSocketContext";


const defaultProfile =

{
	anamorphic
		:
		"off",
	averageBitrate
		:
		15000,
	chromaSmooth
		:
		"off",
	chromaSmoothTune
		:
		"none",
	codec
		:
		"Any",
	codecs
		:
		[],
	color
		:
		"black",
	colorspace
		:
		"off",
	constantQuality
		:
		22,
	container
		:
		"matroska",
	cropping
		:
		"off",
	deblock
		:
		"off",
	deblockTune
		:
		"medium",
	deinterlace
		:
		"off",
	deinterlacePreset
		:
		"default",
	denoise
		:
		"off",
	denoisePreset
		:
		"light",
	denoiseTune
		:
		"none",
	detelecine
		:
		"off",
	encoder
		:
		"",
	extension
		:
		"mkv",
	fastDecode
		:
		false,
	fill
		:
		"none",
	flipping
		:
		false,
	framerate
		:
		"same as source",
	framerateType
		:
		"peak Framerate",
	grayscale
		:
		false,
	id
		:
		1,
	interlaceDetection
		:
		"off",
	level
		:
		"auto",
	limit
		:
		"none",
	mapUntaggedAudioTracks
		:
		true,
	mapUntaggedSubtitleTracks
		:
		true,
	multipassEncoding
		:
		false,
	name
		:
		"Any",
	passThruCommonMetadata
		:
		true,
	preset
		:
		"medium",
	profile
		:
		"auto",
	profileAudioLanguages
		:
		[],
	profileSubtitleLanguages
		:
		[],
	qualityType
		:
		"constant quality",
	rotation
		:
		0,
	sharpen
		:
		"off",
	sharpenPreset
		:
		"medium",
	sharpenTune
		:
		"",
	tune
		:
		"none",
}



const Profiles = () => {
	const wsContext = useContext(WebSocketContext);
	const profiles = wsContext?.data?.profiles;
	const [selectedProfile, setSelectedProfile] = useState<any>({});
	const [isModalOpen, setIsModalOpen] = useState(false);
	const handleProfileClick = (profile: any) => {

		setSelectedProfile(profile);
		setContent(profile);
		setIsModalOpen(true);
	};

	const onModalDelete = async () => {
		await fetch(`/api/profiles/${selectedProfile?.id}`, {
			method: "DELETE",
			headers: {
				"Content-Type": "application/json",
				Authorization: `Bearer ${localStorage.getItem("token")}`,
			},
		});
		setIsModalOpen(false);
	};

	const onModalSave = async () => {
		await fetch(`/api/profiles/` + content.id, {
			method: "PUT",
			headers: {
				"Content-Type": "application/json",
				Authorization: `Bearer ${localStorage.getItem("token")}`,
			},
			body: JSON.stringify(content),
		});
		setIsModalOpen(false);
	};
	const [content, setContent] = useState(selectedProfile);


	return (
		<div className={styles.profiles}>
			<ProfileModal
				header={"Edit - Profile"}
				type={"profile"}
				isOpen={isModalOpen}
				setIsOpen={setIsModalOpen}
				onSave={onModalSave}
				onDelete={onModalDelete}
				data={selectedProfile}
				content={content}
				setContent={setContent}
			/>
			<div className={styles.content}>
				<div className={styles.codecProfiles}>
					<div className={styles.header}>Profiles</div>
					<div className={styles.profileContainer}>
						{profiles?.map((profile: any) => (
							<Profile
								name={profile?.name}
								key={profile?.name}
								codecs={profile?.codecs}
								onClick={handleProfileClick}
								profile={profile}
							/>
						))}
						<Profile
							type={"add"}
							key={"add"}
							name={""}
							codecs={defaultProfile?.codecs}
							profile={{
								...defaultProfile,
								id: profiles?.length + 1,
								name: "",
								codecs: [],
								profileAudioLanguages: [],
								profileSubtitleLanguages: []
							}}
							onClick={handleProfileClick}
						/>
					</div>
				</div>
			</div>
		</div>
	);
};
export default Profiles;
