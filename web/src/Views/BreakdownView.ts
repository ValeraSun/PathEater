import * as THREE from "three";

export function CreateHole(): THREE.Mesh
{
	const geometry = new THREE.CircleGeometry(1.2, 24);
	const material = new THREE.MeshStandardMaterial({
		color: 0x050505,
		side: THREE.DoubleSide,
		roughness: 1,
		metalness: 0
	});

	const mesh = new THREE.Mesh(geometry, material);
	mesh.rotation.x = -Math.PI / 2; // если дыра плоская на полу/стене — подгони ориентацию под сцену

	return mesh;
}